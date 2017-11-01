This is configurable through environment variables `SLAB_ZENDESK_URL` and `SLAB_SLACK_API`. 

The goal is to have the SLAB, or SLA Bot, grab a list of the latest tickets, parse the returned JSON, identify any tickets with SLA policies applied to them, determine the next breach time and send a Slack notification at given intervals for reminders. 

To get the list of tickets, the Zendesk API will be used:

```
/api/v2/tickets.json?include=slas
```

From there, it returns a JSON blob that looks like:

```
{"url":"https://circleci1504710777.zendesk.com/api/v2/tickets/24.json","id":24,"external_id":null,"via":{"channel":"web","source":{"from":{},"to":{},"rel":"web_widget"}},"created_at":"2017-10-25T15:25:50Z","updated_at":"2017-10-25T15:59:41Z","type":null,"subject":"Test SLA","raw_subject":"Test SLA","description":"SLA Bot grab\n\n------------------\nSubmitted from: http://localhost:8080/","priority":"low","status":"open","recipient":null,"requester_id":25755546288,"submitter_id":25755546288,"assignee_id":25670682267,"organization_id":27631063347,"group_id":43489307,"collaborator_ids":[],"follower_ids":[],"forum_topic_id":null,"problem_id":null,"has_incidents":false,"is_public":true,"due_at":null,"tags":["2_0","cloud","not_mobile","platinum","web_widget"],"custom_fields":[{"id":80799627,"value":null},{"id":80799647,"value":"tylerconlee"},{"id":81186648,"value":"not_mobile"},{"id":80799747,"value":"test"},{"id":81186668,"value":"2_0"},{"id":81204108,"value":"cloud"},{"id":81204228,"value":null},{"id":80961507,"value":null}],"satisfaction_rating":{"score":"unoffered"},"sharing_agreement_ids":[],"fields":[{"id":80799627,"value":null},{"id":80799647,"value":"tylerconlee"},{"id":81186648,"value":"not_mobile"},{"id":80799747,"value":"test"},{"id":81186668,"value":"2_0"},{"id":81204108,"value":"cloud"},{"id":81204228,"value":null},{"id":80961507,"value":null}],"ticket_form_id":758767,"brand_id":7706667,"satisfaction_probability":null,"slas":{"policy_metrics":[{"breach_at":"2017-10-26T15:59:41Z","stage":"active","metric":"next_reply_time","hours":24}]},"allow_channelback":false}],"next_page":null,"previous_page":null,"count":9}%
```

The crucial part of this is near the end:

```
"slas":{"policy_metrics":[{"breach_at":"2017-10-26T15:59:41Z","stage":"active","metric":"next_reply_time","hours":24}]}
```

There are 2 "stages": `active` and `completed`. The `next_reply_time` can be in either `hours` or `minutes`. The timezone for `breach_at` appears to be PST, however, this is unconfirmed. 


Within Slack, I'll set up a new channel #support_alerts. The notification scheme would look like:

if `hours` < 3, show yellow color on Slack attachment

if `minutes`, which would mean that the breach is wihtin 1 hour, show red color on Slack attachment

if `minutes` < 30, send an @here to the channel with a red color on Slack attachment