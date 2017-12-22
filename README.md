# slab - Zendesk SLA Bot for Slack

[![CircleCI](https://circleci.com/gh/TylerConlee/slab.svg?style=svg)](https://circleci.com/gh/TylerConlee/slab)
[![GoDoc](https://godoc.org/github.com/TylerConlee/slab?status.svg)](https://godoc.org/github.com/TylerConlee/slab)

> This bot is a Go app that monitors a Zendesk instance and reports upcoming SLA breaches to a given Slack channel. 

![Slab in action](https://user-images.githubusercontent.com/3723686/34063510-670880a2-e1a7-11e7-8f18-7b83afaab60f.gif)

## Installation
Clone the repo and run `go get -t -d -v ./...` to ensure any and all dependencies are local, followed by `go build .`. This creates a binary called `slab` in your local folder. 

## Usage
To run SLAB, create a configuration `.toml` file based off of the [configuration options](https://github.com/TylerConlee/slab/wiki/Configuring-SLAB). The `.toml` file path is then passed as an argument when starting SLAB:

```
./slab config.toml
```

## Documentation

[Full SLAB Documentation](https://github.com/TylerConlee/slab/wiki)

## Contribute

PRs accepted. 

## License
SLAB is released under the [MIT License](https://github.com/TylerConlee/slab/blob/master/LICENSE.md)