<p align="center">
  <img src="https://s3.amazonaws.com/forwardlytics-assets/logo-color.svg">
</p>

Takes event in and forwards them to various places.

You can send events to Forwardlytics, and it will forward it to many services like [Intercom][intercom] or [Mixpanel][mixpanel].

**THIS IS A VERY VERY EARLY NON WORKING VERSION**. Use at your own risk, or contribute :)

Where are we at? What needs to be done and what's in the pipeline? See [Forawrdlytics on waffle.io][forwardlytics-on-waffle]

[![Build Status](https://travis-ci.org/jipiboily/forwardlytics.svg?branch=master)](https://travis-ci.org/jipiboily/forwardlytics)
[![Stories in Progress](https://badge.waffle.io/jipiboily/forwardlytics.svg?label=In%20Progress&title=In%20Progress)](http://waffle.io/jipiboily/forwardlytics)

## Installation

- `go get github.com/jipiboily/forwardlytics`

- set `FORWARDLYTICS_API_KEY=SOMETHING_YOU_DECIDE_AND_NO_ONE_CAN_GUESS`

To send to Intercom:
- set `INTERCOM_API_KEY=123`
- set `INTERCOM_APP_ID=456`

To send to Mixpanel:
- set `MIXPANEL_TOKEN=abc`
- set `MIXPANEL_API_KEY=123`

[Drip][drip] and [Keen.io][keen.io] are probably going to be next.

## Deployment

Forwardlytics can be deployed to [Heroku][heroku]. You can setup the port it starts on by setting the `PORT` environment variable.

## You need an integration that doesn't exist yet?

You have two options:

- send a PR adding it.
- get in touch to have it added by the author (for a fee).

## Calling the API

cURL example:

```
curl --request POST \
--header "Content-Type: application/json" \
--header "Forwardlytics-Api-Key: 123ma" \
-d '{"userID":"123", "userTraits":{"number_of_things":"42"},"timestamp":1459532831}' http://localhost:3000/identify
```

See [./integration/integration.go][integration.go] for details of what is accepted by the API.

## Development

Note that you should install [Godep][godep] if you are going to add any dependency to this project.

- `FORWARDLYTICS_API_KEY=somevalue go run main.go`

If you want auto reloading, install [codegangsta/gin][codegangsta/gin] and run:

- `FORWARDLYTICS_API_KEY=somevalue gin -a 8080 -t .`

You need to set environment variables for the integrations you want to work with.

## Why?

[Read "Self-hosted Segment equivalent?" on Medium][self-hosted-segment-equivalent]

[intercom]: https://www.intercom.io/
[mixpanel]: https://mixpanel.com/
[heroku]: https://www.heroku.com/
[forwardlytics-on-waffle]: https://waffle.io/jipiboily/forwardlytics
[integration.go]: https://github.com/jipiboily/forwardlytics/blob/master/integrations/integration.go
[codegangsta/gin]: https://github.com/codegangsta/gin
[https://github.com/tools/godep]: https://github.com/tools/godep
[self-hosted-segment-equivalent]: https://medium.com/@jipiboily/self-hosted-segment-equivalent-c81815e963df

## Thanks!

Thanks to my friend <a href="https://twitter.com/juliandoesstuff" target="_blank">Julian</a> for the logo! :grinning: