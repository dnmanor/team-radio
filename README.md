# Team Radio

Simple bot to add music to team playlists via slack channels

## Setup

- Set up your app on [spotify](https://developer.spotify.com/documentation/web-api/tutorials/getting-started)
- Set up your bot on [slack](https://api.slack.com/start/building)
  Add the right permissions needed for the bot. Enough to read slack messages and respond to users. Also make sure to have Socket mode on. (note: any changes you make to bot settings wont reflect on slack until you reinstall the app in workspace)
- Add the needed environment variables, find example in `.sample.env`.

## Run

`go run main.go`
