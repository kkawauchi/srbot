#srbot
##prerequisites
* [Go](https://golang.org)
* [reddit dev prefs](https://www.reddit.com/prefs/apps)

##build
```bash
go get -u github.com/uub/srbot
```
##create wiki pages

####Sticky
https://www.reddit.com/r/SUBREDDIT_NAME/wiki/bot/sticky/PAGE_NAME

PAGE_NAME | format | description
--- | --- | ----
title | Text | タイトルの/r/と日付の間
desc | Markdown | 本文の/r/直後
linkstr | Text | 本文の過去リンク(flairでの検索)のテキスト
footer | Markdown | 本文の最後
flair | Text | ポストに追加されるflair(同一のflairがsubreddit内に作成されていなければならない)
interval | Text | 次から選択、day、week、month

##run and post sticky
```bash
srbot
curl http://localhost:8080/sticky?uid=REDDIT_ID&upw=REDDIT_PASSWORD&did=REDDIT_DEV_ID&dpw=REDDIT_DEV_PASSWORD&sr=SUBREDDIT_NAME
```
##for App Engine
create app.yaml, and cron.yaml if you'd like

app.yaml:
```yaml
application: YOUR_APP_ID
version: 1
runtime: go
api_version: go1

handlers:
- url: /sticky*
  script: _go_app
  login: admin
- url: /.*
  script: _go_app
  secure: always
```

cron.yaml:
```yaml
cron:
- description: SOME_DESCRIPTION
  url: /sticky?uid=REDDIT_ID&upw=REDDIT_PASSWORD&did=REDDIT_DEV_ID&dpw=REDDIT_DEV_PASSWORD&sr=SUBREDDIT_NAME
  schedule: every monday 00:00
  timezone: Japan
```
then deploy
