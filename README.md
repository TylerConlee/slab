# slab - Zendesk SLA Bot for Slack

[![CircleCI](https://circleci.com/gh/TylerConlee/slab.svg?style=svg)](https://circleci.com/gh/TylerConlee/slab)
[![GoDoc](https://godoc.org/github.com/TylerConlee/slab?status.svg)](https://godoc.org/github.com/TylerConlee/slab)

> This bot is a Go app that monitors a Zendesk instance and reports upcoming SLA breaches to a given Slack channel.

![Slab in action](https://user-images.githubusercontent.com/3723686/34063510-670880a2-e1a7-11e7-8f18-7b83afaab60f.gif)

## Installation

Before deployment, a Slack app has to be created in your Slack team. An app can be created [through Slack's web UI](https://api.slack.com/apps).
The Slab application requires several configuration options set in the Slack interface.
Interactive Components must be enabled, with a SSL-enabled link to a server that Slab will run on, ending in `/slack`.

![interactive-components](https://user-images.githubusercontent.com/3723686/36488544-8d829c8e-16d8-11e8-9bc1-9f9a2ec403ed.png)

A Bot User is also required with the username `slab`:

![bot-user](https://user-images.githubusercontent.com/3723686/36488590-ae53968e-16d8-11e8-9b69-19e3c7c1f451.png)

Finally, the OAuth Tokens & Redirect URLs page will provide the Bot User OAuth token, which is used in your configuration `.toml`.
To run SLAB on your server, create a configuration `.toml` file based off of the [configuration options](https://github.com/TylerConlee/slab/wiki/Configuring-SLAB).
The `.toml` file path is then passed as an argument when starting SLAB:

```sh
./slab config.toml
```

## Development setup

Glide must be installed to compile Slab. Clone the repo and run `glide install` to ensure any and all dependencies are local.

```sh
glide install
```

## Release History

* v1.2
  * **[FEATURE]** - More Info - When a ticket sends an SLA notification to Slack,
    you can now see more information about the ticket and who requested it with the new More Info button
  * **[FEATURE]** - Colored notifications - based off of how long the ticket has left before the SLA timer expires.
  * **[CHANGE]** - Updated logging to utilize sirupsen/logrus
* v1.1
  * **[FEATURE]** - Triager role - enable the notification of new tickets through direct message,
  and stay on top of tickets as they come in to Zendesk
  * **[FEATURE]** - Added Slack commands `@slab set`, `@slab unset`, `@slab whois`, `@slab help` `@slab status`.
    * **@slab set** -
    Set returns a drop down menu of all available Slack members, allowing you to select a person to take the triager role. Note that this is all employees, and not just support. For the most part, this command will be used when setting yourself as triager.

    * **@slab unset** -
    Unset returns Slab to its default state of having no triager. This means that the new ticket notifications will not be sent until a new triager is set.

    * **@slab whois** -
    Whois returns the name of the person currently set as triager. If there is not currently a triager, or if unset was recently ran, this value will be None.

    * **@slab help** -
    Help returns a list of all available Slab commands.

    * **@slab status** -
    Status returns the server's metadata, including what version of Slab is running and what the uptime is.
  * **[BUGFIX]** - Updated fork of Slack dependency to avoid rate limiting issues
* v1.0
  * **[FEATURE]** - Configurable port settings
  * **[FEATURE]** - Acknowledge button added for SLA breach messages
  * **[BUGFIX]** - SLAB crashes when EOF reached in Zendesk API

## Meta

Tyler Conlee – [@TylerConlee](https://twitter.com/tylerconlee) – tyler@circleci.com

Distributed under the MIT license. See ``LICENSE`` for more information.

[https://github.com/tylerconlee/](https://github.com/dbader/)

## Contributing

1. Fork it (<https://github.com/yourname/yourproject/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

[wiki](https://github.com/yourname/yourproject/wiki)
