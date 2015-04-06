#srbot
##prerequisites
* [Go](https://golang.org)
* [reddit dev prefs](https://www.reddit.com/prefs/apps)

##build
```bash
go get -u github.com/uub/srbot
```
##create wiki pages
create wiki pages on subreddit respectively by the feature

###sticky
https://www.reddit.com/r/SUBREDDIT_NAME/wiki/bot/sticky/PAGE_NAME

PAGE_NAME | format | description
--- | --- | ----
title | Text | between /r/ and dates in title
desc | Markdown | description in text after /r/
linkstr | Text | label of flair link in text
footer | Markdown | footer in text
flair | Text | fair text existed already in a list of subreddit flairs
interval | Text | choose from "day", "week" and "month" (not included quotation)

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

example of cron.yaml that works once a week:
```yaml
cron:
- description: SOME_DESCRIPTION
  url: /sticky?uid=REDDIT_ID&upw=REDDIT_PASSWORD&did=REDDIT_DEV_ID&dpw=REDDIT_DEV_PASSWORD&sr=SUBREDDIT_NAME
  schedule: every monday 00:00
  timezone: Japan
```
then deploy
