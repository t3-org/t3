- Create channels in the config.

```yaml
channels:
  matrix: # the map's key is the channel name.
    - type: matrix
      homeserver: "https://..."
      username: "..."
      password: "..."


```

- Create webhooks. Each webhook can have a channel optionally (spcify the channel by its name).
- Grafana will send a firing alert to our webhook.
- Create the ticket (if it's not already created by another webhook, check by its fingerprint). and then
  get its channel instance (we'll have a map of channels instances) and calls it's `firing` method to send
  a message to the target third-party service(e.g., matrix).
- Do the same for resolved grafana alert. I mean update its state in DB(if it's not upated already) and use
  the webhook's channel to send a message.

- Now we should let various interfaces update the ticket data(matrix commands, ui,cmd...).
- I'll describe for matrix.
- Per each channel call its `start` method which should start listening to messages. just like the api server that
  we call its listen at the start time of the app, do the same per each channel instance to listen to the messages.
- for matrix follow the matrix specs to see how should be the message format to do some action.

### Matrix messages format.

Matrix commands have a configured prefix(default: `!!`). we'll use `!!` as our prefix in the next sections of the doc.

__Matrix commands__

- `!!seen {minutes(default: 0)}`. set the `seen_at` of a ticket. e.g. `!!seen -10` meaning i've seen 10 minutes ago.
- `!!spam`, `!!spam true` `!!spam false`. set the spam field of a ticket.
- `!!resolved {minutes(default: 0)}`. e.g., `!!resolved` the alert is resolved now.
- `!!firing`. just changes the `is_firing` field to `true`.
- `!!level {level: low,medium or high}`: the set level of a ticket. e.g., `!!level low`.
- `!!description {msg}`. set the `description` field(it'll remove the previous content if it's not empty).
- `!!init_description {msg}`. set the `description` field if it's empty.

- `!!edit`: Get a link to the ticket edition form in the UI.
- `!!ticket`: Get the ticket.

Non-ticket-related commands (commands that are not for specific command)

- `!!new`: get a link to the ticket creation form in the UI.
- `!!dashboard` get a link to the dashboard.
- `!!help` get a help message for all commands and also some useful links like dashboard link.
