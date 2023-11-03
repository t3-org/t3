tickets:

- *id(int)
- *fingerprint(string)

- *status (enum: firing,resolved)
- *is_spam(bool)
- level(enum: low,medium,high):

- started_at(datetime)
- reported_at? (sre report time, I'm not sure if we need to this field!)
- seen_at(datetime)
- ended_at(datetime)

- description (text, fill-in manually)

ticket_tags:

- id(int)
- ticket_id(int)
- term(string): some expected tags:
    - "oncall:mehran_prs" (we should get this automatically)
    - "sre_oncall:reza" (we should get this automatically)
    - "team:routing" (set on grafana label)
    - "zone:teh1" (set on grafana label)
    - "service:air" (set as grafana label)
    - "cause:network", "cause:empty_cassandra"
    - "action:report_to_cloud","action:patch_service_and_deploy"

webhooks:

- id
- channel_type(string): type of the channel that we want to send message to it. e.g., `matrix`.
- channel_config(string): config of the channel, for matrix it's the room_id, for other channels could be `json`.

tickets_attached_entities:

- ticket_id(int)
- entity_id(int): a unique id of an entities that attaches to a ticket. e.g., an matrix event.
  use it to find the ticket which is attached to a matrix's event.