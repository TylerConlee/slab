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

## Development schedule


## Meta

Tyler Conlee – [@TylerConlee](https://twitter.com/tylerconlee) – tyler@circleci.com

Distributed under the GPU license. See ``LICENSE`` for more information.

[https://github.com/tylerconlee/](https://github.com/tylerconlee/)

## Contributing

1. Fork it (<https://github.com/tylerconlee/slab/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

