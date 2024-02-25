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

ticket_labels:

- ticket_id(int)
- key (text)
- val (text)
