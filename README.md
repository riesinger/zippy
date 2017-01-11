# :zap: zippy

Is a small little CMS written in Go. I use it for my personal homepage and thought I might share it.

## Features
- Online Markdown Editor for writing articles
- Theming support with Handlebars
- Admin Backend (WIP)
- Stores Articles in MongoDB (as I already have it running on my VSP and do not want to clutter it. Other DB's could possibly be implemented)

#### Planned
- Authentication (for multiple authors)
- More page variables for Handlebars templates
- More standard themes
- Upload articles as files (for writing offline)
- Upload themes in admin panel
- See some stats in admin panel

## Install
Clone this repo into your `$GOPATH` and then run `go install github.com/arial7/zippy`. You currently need to run the server
from its src dir!

```bash
mkdir -p $GOPATH/src/github.com/arial7/
git clone https://github.com/arial7/zippy $GOPATH/src/github.com/arial7/
go get github.com/arial7/zippy
go install github.com/arial7/zippy
```

This install process is quite tedious, I know... I will make it better soon :tm:!

## Running
Make sure, you have `$GOPATH/bin` in your `$PATH`, then just run `zippy` from its source dir.
