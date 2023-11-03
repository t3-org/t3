**Firing**

Labels:

- alertname = T2
- team = blue
- zone = us-1
  Annotations:
- description = This is the alert rule checking the second system
- runbook_url = https://myrunbook.com
- summary = This is my summary
  Source: https://play.grafana.org/alerting/1afz29v7z/edit
  Silence: https://play.grafana.org/alerting/silence/new?alertmanager=grafana&matchers=alertname%3DT2%2Cteam%3Dblue%2Czone%3Dus-1
