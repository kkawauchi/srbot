// Package bot provides subreddit bot
package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Context struct {
	Client *http.Client
	W      http.ResponseWriter
	R      *http.Request

	token string

	// By url params
	dev_id            string
	dev_pass          string
	user_id           string
	user_pass         string
	subreddit         string
	api_flairselector apiUrl
	api_selectflair   apiUrl

	// By wiki (/r/SUBREDDIT/wiki/bot/sticky/PAGE)
	title        string
	text_desc    string
	text_linkstr string
	text_footer  string
	interval     postInterval
	flair        string
}

const text = `
/r/%s%s

[%s](http://www.reddit.com/r/%s/search?sort=new&restrict_sr=on&q=flair:%s)

%s
`

type postInterval int

const (
	post_per_day postInterval = iota
	post_per_week
	post_per_month
)

type apiUrl string

const (
	api_access_token         apiUrl = `https://www.reddit.com/api/v1/access_token`
	api_submit               apiUrl = `https://oauth.reddit.com/api/submit`
	api_set_subreddit_sticky apiUrl = `https://oauth.reddit.com/api/set_subreddit_sticky`
	api_r                           = `https://oauth.reddit.com/r/`
)

type ResToken struct {
	Access_token string `json:"access_token"`
	Token_type   string
	Expires_in   int
	scope        string
}

type ResJson struct {
	Kind string `json:"kind"`
	JSON struct {
		Errors []interface{} `json:"errors"`
		DATA   struct {
			Url  string
			Id   string
			Name string `json:"name"`
		}
	}
	Current ResFlair
	Choices []ResFlair
	DATA    struct {
		Children   []ResText
		Content_md string `json:"content_md"`
	}
}
type ResFlair struct {
	Flair_css_class     string `json:"flair_css_class"`
	Flair_template_id   string `json:"flair_template_id`
	Flair_text_editable bool
	Flair_position      interface{}
	Flair_text          string `json:"flair_text"`
}
type ResText struct {
	DATA struct {
		Selftext      string `json:"selftext"`
		Selftext_html string `json:"selftext_html"`
	}
}

type ReqSubmit struct {
	title string
	text  string
}

func init() {
	http.HandleFunc("/sticky", func(w http.ResponseWriter, r *http.Request) {
		c := newContext(w, r)

		uid := r.FormValue("uid")
		upw := r.FormValue("upw")
		did := r.FormValue("did")
		dpw := r.FormValue("dpw")
		sr := r.FormValue("sr")
		if uid == "" || upw == "" || did == "" || dpw == "" || sr == "" {
			fmt.Fprintln(w, "need reddit uid(user id), upw(user password), did(developer id), dpw(developer password), sr(subreddit) for params")
			return
		}

		c.user_id = uid
		c.user_pass = upw
		c.dev_id = did
		c.dev_pass = dpw
		c.subreddit = sr

		c.api_flairselector = apiUrl(api_r + c.subreddit + `/api/flairselector`)
		c.api_selectflair = apiUrl(api_r + c.subreddit + `/api/selectflair`)

		if err := c.sticky(); err != nil {
			fmt.Fprintln(w, err.Error())
			return
		} else {
			fmt.Fprintln(w, "done")
			return
		}
	})
}

func (c *Context) sticky() error {
	var err error
	var j ResJson

	// Get token
	c.logDebug("GET TOKEN\n")

	data := url.Values{
		"grant_type": {"password"},
		"username":   {c.user_id},
		"password":   {c.user_pass},
	}.Encode()

	var t ResToken
	err = c.post(api_access_token, data, &t)
	if err != nil {
		return err
	}

	c.token = t.Access_token

	// Get data from wiki
	c.logDebug("GET DATA FROM WIKI\n")

	err = c.get(api_r+c.subreddit+"/wiki/bot/sticky/title", &j)
	if err != nil {
		return err
	}
	c.title = j.DATA.Content_md

	err = c.get(api_r+c.subreddit+"/wiki/bot/sticky/desc", &j)
	if err != nil {
		return err
	}
	c.text_desc = j.DATA.Content_md

	err = c.get(api_r+c.subreddit+"/wiki/bot/sticky/linkstr", &j)
	if err != nil {
		return err
	}
	c.text_linkstr = j.DATA.Content_md

	err = c.get(api_r+c.subreddit+"/wiki/bot/sticky/footer", &j)
	if err != nil {
		return err
	}
	c.text_footer = j.DATA.Content_md

	err = c.get(api_r+c.subreddit+"/wiki/bot/sticky/flair", &j)
	if err != nil {
		return err
	}
	c.flair = j.DATA.Content_md

	err = c.get(api_r+c.subreddit+"/wiki/bot/sticky/interval", &j)
	if err != nil {
		return err
	}
	switch j.DATA.Content_md {
	case "month":
		c.interval = post_per_month
	case "week":
		c.interval = post_per_week
	default:
		c.interval = post_per_day
	}

	// Submit
	c.logDebug("SUBMIT\n")

	var s ReqSubmit
	err = c.submit(&s)
	if err != nil {
		return err
	}

	data = url.Values{
		"api_type": {"json"},
		"kind":     {"self"},
		"sr":       {c.subreddit},
		"text":     {s.text},
		"title":    {s.title},
	}.Encode()

	err = c.post(api_submit, data, &j)
	if err != nil {
		return err
	}

	id := j.JSON.DATA.Name

	// Stick it on top
	c.logDebug("STICK IT ON TOP\n")

	data = url.Values{
		"api_type": {"json"},
		"id":       {id},
		"state":    {"true"},
	}.Encode()

	err = c.post(api_set_subreddit_sticky, data, &j)
	if err != nil {
		return err
	}

	// Get flair
	c.logDebug("GET FLAIR\n")

	data = url.Values{
		"api_type": {"json"},
		"link":     {id},
	}.Encode()

	err = c.post(c.api_flairselector, data, &j)
	if err != nil {
		return err
	}

	var fId string
	for _, f := range j.Choices {
		if f.Flair_text == c.flair {
			fId = f.Flair_template_id
			break
		}
	}
	if fId == "" {
		return fmt.Errorf("No %s flair in your subreddit", c.flair)
	}

	// Select flair
	c.logDebug("SELECT FLAIR\n")

	data = url.Values{
		"api_type":          {"json"},
		"flair_template_id": {fId},
		"link":              {id},
	}.Encode()

	err = c.post(c.api_selectflair, data, &j)
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) submit(v *ReqSubmit) error {
	now := time.Now()
	if now.Hour() > 12 {
		now = now.Add(time.Hour * 12)
	}
	y, m, d := now.Date()
	switch c.interval {
	case post_per_month:
		v.title = fmt.Sprintf("/r/%s %s %d年%d月", c.subreddit, c.title, y, int(m))
	case post_per_week:
		now = now.Add(time.Hour * 144)
		_, m2, d2 := now.Date()
		v.title = fmt.Sprintf("/r/%s %s %d年%d月%d日〜%d月%d日", c.subreddit, c.title, y, int(m), d, int(m2), d2)
	default:
		v.title = fmt.Sprintf("/r/%s %s %d年%d月%d日", c.subreddit, c.title, y, int(m), d)
	}

	v.text = fmt.Sprintf(text, c.subreddit, c.text_desc, c.text_linkstr, c.subreddit, c.flair, c.text_footer)
	return nil
}

func (c *Context) get(url string, v interface{}) error {
	var err error

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "bearer "+c.token)
	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//c.logDebug(fmt.Sprintln("RAW:", string(b)))
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil

}

func (c *Context) post(url apiUrl, data string, v interface{}) error {
	var err error

	req, err := http.NewRequest("POST", string(url), bytes.NewBufferString(data))
	if err != nil {
		return err
	}
	if url == api_access_token {
		req.SetBasicAuth(c.dev_id, c.dev_pass)
	} else {
		req.Header.Set("Authorization", "bearer "+c.token)
	}
	req.Header.Set("User-Agent", "appengine:com.appspot.gamedevja:v1.0.0 (by /u/leuro)")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	c.logDebug(fmt.Sprintln("RAW:", string(b)))
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}
