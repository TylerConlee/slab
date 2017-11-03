# slab - CircleCI SLA Bot for Zendesk

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
[Zendesk]
    User = "YOUR-ZENDESK-USERNAME"
    APIKey = "YOUR-ZENDESK-API-KEY"
    URL = "URL-TO-YOUR-ZENDESK-INSTANCE"
[SLA]
    [SLA.LevelOne]
        Low = "8h"
        Normal = "4h"
        High = "1h"
        Urgent = "59m"
    [SLA.LevelTwo]
        Low = "8h"
        Normal = "4h"
        High = "1h"
        Urgent = "59m"
    [SLA.LevelThree]
        Low = "8h"
        Normal = "4h"
        High = "1h"
        Urgent = "59m"
    [SLA.LevelFour]
        Low = "8h"
        Normal = "4h"
        High = "1h"
        Urgent = "59m"
[Slack]
    User = "YOUR-SLACK-USERNAME"
    APIKey = "YOUR-SLACK-API-KEY"
    URL = "URL-TO-YOUR-SLACK-INSTANCE"
```

Zendesk user and API key can be found within the Settings of your Zendesk instance. More information can be found [here] on generating a Zendesk API key and username. SLA times are in the form of `8h`, `9m`, `1s`, etc. Up to 4 different SLA policies can be active at a given time. The Slack integration is done as part of a custom Slack integration that you must set up separately from this bot. Once that has been set up, your API information is used in the configuration to post messages to a given channel.