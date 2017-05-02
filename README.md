<p align="center">
  <img src="https://s3.amazonaws.com/forwardlytics-assets/logo-color.svg">
</p>

Forwards analytics events and identification to various marketing & analytics platforms.

You can send events to Forwardlytics, and it will forward it to the configured services like [Intercom][intercom] or [Drip][drip].

Looking for a Forwardlytics client library? There is one for [Ruby just here](https://github.com/jipiboily/forwardlytics-ruby).

## Status

Where are we at? Can you use it in production?

It is used by [Metrics Watch][metricswatch] for now and it has been for a while now. I would not say it is rock solid but it has been working well enough for us.

What needs to be done and what's in the pipeline? See [Forwardlytics on waffle.io][forwardlytics-on-waffle]

[![Build Status](https://travis-ci.org/jipiboily/forwardlytics.svg?branch=master)](https://travis-ci.org/jipiboily/forwardlytics)
[![Stories in Progress](https://badge.waffle.io/jipiboily/forwardlytics.svg?label=In%20Progress&title=In%20Progress)](http://waffle.io/jipiboily/forwardlytics)

## Installation

- `go get github.com/jipiboily/forwardlytics`

- set `FORWARDLYTICS_API_KEY=SOMETHING_YOU_DECIDE_AND_NO_ONE_CAN_GUESS`

To send to [Intercom][intercom]:
- set `INTERCOM_API_KEY=123`
- set `INTERCOM_APP_ID=456`

To send to [Drip][drip]:

- set `DRIP_ACCOUNT_ID=234` (found here: https://www.getdrip.com/{drip_account_id}/settings/site under "3rd party integrations")
- set `DRIP_API_TOKEN=432` (found here: https://www.getdrip.com/user/edit under "API-token")

**Please note** that you need to send an "email" property to be able to get the Drip integration working.

To send to Drift:

- set `DRIFT_ORG_ID=456` (ATM only possible to find by contacting the drift support dept)

To send to[Mixpanel][mixpanel]:
- set `MIXPANEL_TOKEN=123`

## Deployment

Forwardlytics can be deployed to [Heroku][heroku]. You can setup the port it starts on by setting the `PORT` environment variable.

## Error tracking

Right now Forwardlytics supports tracking error via Bugsnag. Thanks to Logrus, it's pretty easy to add any other bug tracker. PRs welcome.

## Retrying calls on failure

Forwardlytics has a built-in retry-mechanism than can be enabled
should calls to a provider fail. To enable this, set the environment
variable `NUM_RETRIES_ON_ERROR=X` where `X` is the number of retries
to attempt before giving up. This is implemented as an
[exponential backoff algorithm](https://en.wikipedia.org/wiki/Exponential_backoff).


### Bugsnag config

To enable Bugsnag, set those environment variables:

```
BUGSNAG_API_KEY=your-api-key-123
ENVIRONMENT=development
```

If the environment is not set, it'll work but defaults to `development`.

## You need an integration that doesn't exist yet?

You have two options:

- send a PR adding it.
- [get in touch to have it added by the author (for a fee)][email].

## How to add a new integration

To add a new integration you need to add a package that implements the
[Integration interface](integrations/integration.go) to a separate
folder in the [integrations/](integrations/) subfolder of this
project, usually named after the integration. The integration should
be toggled by adding an ENV-variable that is picked up by the
`Enabled()`-function in the integration and that is passed to
forwardlytics on startup. To activate the new integration, add the
path to the new integration in the import-statement in
[main.go](main.go). Remember to add an `init()` function to the new
package that registers the new integration using
`integrations.RegisterIntegration(<integration-name>,
integration)`. For examples, see the different integrations in the
[integrations/](integrations/) subfolder
(eg. [the drip-integration](integrations/drip/drip.go)). Don't forget
to add tests for all endpoints and for other integration spesific
stuff.

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

## Thanks!

Thanks to my friend <a href="https://twitter.com/juliandoesstuff" target="_blank">Julian</a> for the logo! :grinning:

[email]: mailto:jp@metrics.watch
[metricswatch]: http://metricswatch.com
[intercom]: https://www.intercom.io/
[mixpanel]: https://mixpanel.com/
[drip]: http://getdrip.com/
[keen.io]: http://keen.io/
[heroku]: https://www.heroku.com/
[forwardlytics-on-waffle]: https://waffle.io/jipiboily/forwardlytics
[integration.go]: https://github.com/jipiboily/forwardlytics/blob/master/integrations/integration.go
[codegangsta/gin]: https://github.com/codegangsta/gin
[https://github.com/tools/godep]: https://github.com/tools/godep
[self-hosted-segment-equivalent]: https://medium.com/@jipiboily/self-hosted-segment-equivalent-c81815e963df
