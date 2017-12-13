# slab - Zendesk SLA Bot for Slack


[![CircleCI](https://circleci.com/gh/TylerConlee/slab.svg?style=svg)](https://circleci.com/gh/TylerConlee/slab)

This bot is a Go app that monitors a Zendesk instance and reports upcoming SLA breaches to a given Slack channel. 

## Installation
To install SLAB, add a configuration `.toml` file to the directory that the compiled binary is located in. The `.toml` file path is then passed as an argument when starting SLAB:

```
./slab config.toml
```

## Configuration
The configuration `.toml` contains everything that SLAB needs in order to know where to get the tickets from, and where to send the tickets to, as well as what the SLA structure looks like. This allows for complete customization of SLA policies and means that SLAB is flexible enough for use with any Zendesk or Slack instance. 

Within your configuration file, you should have this structure:

```
UpdateFreq = "10m"
LogLevel = "Debug"
[Zendesk]
    User = "YOUR-ZENDESK-USERNAME"
    APIKey = "YOUR-ZENDESK-API-KEY"
    URL = "URL-TO-YOUR-ZENDESK-INSTANCE"
[Slack]
    APIKey="YOUR-API-KEY"
    ChannelID="TARGET-CHANNEL-ID"
[SLA]
    [SLA.LevelOne]
        Tag = "platinum"
        Low = "8h"
        Normal = "4h"
        High = "1h"
        Urgent = "59m"
        Notify = true
    [SLA.LevelTwo]
        Tag = "gold"
        Low = "8h"
        Normal = "4h"
        High = "1h"
        Urgent = "59m"
        Notify = true
    [SLA.LevelThree]
        Tag = "silver"
        Low = "8h"
        Normal = "4h"
        High = "1h"
        Urgent = "59m"
        Notify = true
    [SLA.LevelFour]
        Tag = "standard"
        Low = "8h"
        Normal = "4h"
        High = "1h"
        Urgent = "59m"
        Notify = false

```

`UpdateFreq` - in the form of `8h`, `9m`, `1s`, etc. Used to set the update loop. Default should be set to `10m`. 

`LogLevel` - sets the log level output level. `debug` `info` and `notice` are valid options


Zendesk user and API key can be found within the Settings of your Zendesk instance. More information can be found [here](https://support.zendesk.com/hc/en-us/articles/226022787-Generating-a-new-API-token-) on generating a Zendesk API key and username. SLA times are in the form of `8h`, `9m`, `1s`, etc. Up to 4 different SLA policies can be active at a given time. The Slack integration is done as part of a [bot user](https://api.slack.com/bot-users) in Slack. Once that has been set up, you're given a token for your bot. 