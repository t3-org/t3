### Monitoring contact-points vs T3 channels:
Monitoring systems want to inform the team for the alert. But T3 want to have
interactive channels to let the team members interact with the incident manager.  
So this tool should contain the contact points which are interactive(like matrix, slack or Telegram bot).
So if we're using this tool, it doesn't mean we should replace contract-points of setting of let's say
Grafana with T3 channels settings.

So let all you contact-points remain in Grafana(or your monitoring tool) and just add some interactive
channels here.

### How does channel-system work?
We have various channel types (e.g., matrix, slack,...)

We'll connect tickets to channels using tickets labels.

We'll implement it by defining channels, channel_policies in T3 config.

- Channel: it's the config for specific channel. each channel can have a `base` field which specifies its
  base configs. the specified `base` itself can not have another `base` field.
- channel_policies: we'll use them to select channels per each ticket based on the ticket's labels.
  we'll do `AND` between labels. and the specific labels in the config can just be equal to the ticket's
  labels to mark that ticket as matched with the policy's labels.

An example config could be like the following snippet:

```yaml
channel_homes:
  central_matrix: # the server name is central_matrix.
    type: matrix
    config:
      username: abc
      password:
        env: BASE_MATRIX_PASSWORD # it supports reading from env variable too.
      
channels:
  sre:
    home: central_matrix
    config:
      room_id: "!asodifjewfasff"
  orders:
    home: central_matrix
    config:
      room_id: "!3ifjslanfasdfadsvd"

channel_policies:
  - channel: sre # send all messages to the "sre" channel.
  - channel: orders # send tickets with label "team: orders" to the "orders" channel.
    labels:
      team: orders
```

