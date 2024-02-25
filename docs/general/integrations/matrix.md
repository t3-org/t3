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


### Matrix Links
- [matrix concepts](https://spec.matrix.org/v1.8/client-server-api/#sending-events-to-a-room)
- [matrix api](https://spec.matrix.org/v1.8/client-server-api/#sending-events-to-a-room)
- [matrix sdk examples](https://github.com/matrix-org/matrix-rust-sdk/tree/main/examples)
- [matrix element client](https://app.element.io)


### TODO
- Sometimes users send extra params and think that the param is applied, but the command doesn't
  accept any param, in these cases we can check if the command has any param, return error.
