Used global fields:

- `alerts`

ignored global fields:

- `receiver`
- `status`
- `orgId`: We can later use it to support multi-tenancy.
- `truncatedAlerts`
- `groupLabels`
- `commonLabels`
- `commonAnnotations`
- `externalURL`
- `version`
- `groupKey`
- `truncatedAlerts`

Used Alert fields:

- `status`: is firing or not.
- `labels`: the alert labels.
- `annotations`:? the alert annotations.
- `startsAt`
- `endsAt`
- `values`: The alert values that triggered this alert.
- `generatorURL`: url to the alert on grafana panel.
- `fingerprint`: The labels fingerprint, alarms with the same labels will have the same fingerprint.

Ignored alert fields:

- `silenceURL`
- `imageURL`

