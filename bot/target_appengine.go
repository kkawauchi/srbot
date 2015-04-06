// +build appengine

package bot

import (
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	g := appengine.NewContext(r)
	return &Context{
		Client: urlfetch.Client(g),
		W:      w,
		R:      r,
	}
}

func (c *Context) logInfo(s string) {
	g := appengine.NewContext(c.R)
	g.Infof("%s", s)
}

func (c *Context) logDebug(s string) {
	g := appengine.NewContext(c.R)
	g.Debugf("%s", s)
}

func (c *Context) logError(s string) {
	g := appengine.NewContext(c.R)
	g.Errorf("%s", s)
}
