### Prerequisites
- Setup [channels config](./channels.md).
- Register a webhook in Grafana to send tickets to T3.

### Ticket Flow
- Grafana will send a firing alert to our webhook.
- When it gets an alert in the firing state:
  - Creates a ticket for the alert  (with started_at,tags,...) if it's not created already(by another webhook)
  - Send a message to the channels that match with the ticket's labels (by calling the `firing` method of each channel).
  - 

- When it get a alert in the resolved state:
  - Sets the `end_date` of the alert(if it's not set already by the users) by finding the alert in the database
    using `fingerprint` field.


- Now we should let various interfaces update the ticket data(matrix commands, ui,cmd...).

### User actions

- users should set `level`, `seen_at` of each ticket.
- Users can set all ticket's fields if they want.

Users do actions via the following ways:

- cli
- ui
- channels(e.g., through matrix UI).

__notifications that we'll send to the channels(e.g., matrix room)__

- ticket created(firing).
- ticket reoslved





