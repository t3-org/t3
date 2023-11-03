__Flow__

- Our app registers an webhook on grafana to get all alerts.
- When it gets an alert(in the firing state):
    - Creates a ticket for the alert  (with started_at,tags,...) if it's not created already(by another webhook)
    - Send a message to the channel that is assigned to the webhook (e.g., `matrix` channel).

- When it get a alert in the resolved state:
    - Sets the `end_date` of an alert(if it's not set already by the users) by finding the alert in the database
      using `fingerprint` field.

__User actions__

- users should set `level`, `seen_at` of each ticket.
- Users can set all ticket's fields if they want.

Users do actions using the following waus:

- cli
- ui
- matrix commands.

__notifications that we'll send to the channels(matrix room)__

- ticket created(firing).
- ticket reoslved

__automatic actions__
We'll set some extra tags on each ticket automatically by fetching data from some palces like on_call table. the
following fields will be set automatically:

- team's oncall.
- sre_oncall.

We'll detect if a ticket should be marked as spam or not by the following rule:

- if the alert is because of `no_data` or `QueryError` on grafana.

