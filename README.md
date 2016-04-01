# Forwardlytics

Takes event in and forwards them to various places.

You can send events to Forwardlytics, and it will forward it to many services like Intercom or Mixpanel.

**THIS IS A VERY VERY EARLY NON WORKING VERSION**. Use at your own risk, or contribute :)

## Installation

- `go get github.com/jipiboily/forwardlytics`

- set `FORWARDLYTICS_API_KEY=SOMETHING_YOU_DECIDE_AND_NO_ONE_CAN_GUESS`

To send to Intercom:
- set `INTERCOM_API_KEY=123`
- set `INTERCOM_APP_ID=456`

To send to Mixpanel:
- set `MIXPANEL_TOKEN=abc`
- set `MIXPANEL_API_KEY=123`

Drip and Keen.io are probably going to be next.

## Deployment

Make sure that everything goes through SSL.

[TBD]

## Development

- `FORWARDLYTICS_API_KEY=somevalue go run main.go`

If you want auto reloading, install https://github.com/codegangsta/gin and run:

- `FORWARDLYTICS_API_KEY=somevalue gin -a 8080 -t .`

You need to set environment variables for the integrations you want to work with.

## Plan for v0.1

Very first version can support only identification of users, and send them to Intercom and Mixpanel only. Next step would be tracking events in the backend. Everything front-end related is for later.

For now, Forwardlytics supports only one project at a time. The configuration is all in environment variables, no database or anything.

### TODO

- [ ] define and document the API structure
- [ ] implement the Intercom integration
- [ ] add tests for the plugin system, and everything else (why not before? This is a PoC, not sure the structure of the project and data will stay the same yet)
- [ ] implement the Mixpanel integration
- [ ] document how to add an integration
- [ ] implement the Ruby SDK (https://github.com/jipiboily/forwardlytics-ruby)
- [ ] make it robust enough to use in small apps, that are not too critical to start with

## Later...
- add Drip integration
- add Keen.io integration
- have a front-end library that will track things like pageviews and so on
- configure multiple projects, with a database
- store all the same information as analytics.js
- have a web ui
- store all the events in a database
- have a live view of the events coming
- replay events stored in DB when adding an new integration