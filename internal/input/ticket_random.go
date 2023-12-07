package input

import (
	"time"

	"github.com/kamva/gutil"
)

func RandomCreatTicket() *CreateTicket {
	raw := `{
  "receiver": "tmp_webhook",
  "status": "firing",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "Memory Usage",
        "grafana_folder": "my_folder",
        "label_one": "val_one"
      },
      "annotations": {
        "summary": "hichi, yechi kharabe"
      },
      "startsAt": "2023-09-24T17:18:00Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://localhost:3000/alerting/grafana/yeMMX4WSk/view",
      "fingerprint": "3e30172082b00f9a",
      "silenceURL": "http://localhost:3000/alerting/silence/new?alertmanager=grafana\u0026matcher=alertname%3DMemory+Usage\u0026matcher=grafana_folder%3Dmy_folder\u0026matcher=label_one%3Dval_one",
      "dashboardURL": "http://localhost:3000/d/pMEd7m0Mz",
      "panelURL": "http://localhost:3000/d/pMEd7m0Mz?viewPanel=9",
      "valueString": "[ var='B0' metric='Value' labels={} value=7.9525888e+08 ], [ var='B1' metric='cadvisor' labels={name=cadvisor} value=4.325376e+07 ], [ var='B2' metric='grafana' labels={name=grafana} value=4.9606656e+07 ], [ var='B3' metric='prometheus' labels={name=prometheus} value=1.17710848e+08 ], [ var='B4' metric='tlstun' labels={name=tlstun} value=675840 ]"
    }
  ],
  "groupLabels": {
    "alertname": "Memory Usage",
    "grafana_folder": "my_folder"
  },
  "commonLabels": {
    "alertname": "Memory Usage",
    "grafana_folder": "my_folder",
    "label_one": "val_one"
  },
  "commonAnnotations": {
    "summary": "hichi, yechi kharabe"
  },
  "externalURL": "http://localhost:3000/",
  "version": "1",
  "groupKey": "{}/{label_one=\"val_one\"}:{alertname=\"Memory Usage\", grafana_folder=\"my_folder\"}",
  "truncatedAlerts": 0,
  "orgId": 1,
  "title": "[FIRING:1] Memory Usage my_folder (val_one)",
  "state": "alerting",
  "message": "**Firing**\n\nValue: [ var='B0' metric='Value' labels={} value=7.9525888e+08 ], [ var='B1' metric='cadvisor' labels={name=cadvisor} value=4.325376e+07 ], [ var='B2' metric='grafana' labels={name=grafana} value=4.9606656e+07 ], [ var='B3' metric='prometheus' labels={name=prometheus} value=1.17710848e+08 ], [ var='B4' metric='tlstun' labels={name=tlstun} value=675840 ]\nLabels:\n - alertname = Memory Usage\n - grafana_folder = my_folder\n - label_one = val_one\nAnnotations:\n - summary = hichi, yechi kharabe\nSource: http://localhost:3000/alerting/grafana/yeMMX4WSk/view\nSilence: http://localhost:3000/alerting/silence/new?alertmanager=grafana\u0026matcher=alertname%3DMemory+Usage\u0026matcher=grafana_folder%3Dmy_folder\u0026matcher=label_one%3Dval_one\nDashboard: http://localhost:3000/d/pMEd7m0Mz\nPanel: http://localhost:3000/d/pMEd7m0Mz?viewPanel=9\n"
}
`
	severity := "low"
	r := gutil.RandNumber(0, 10)
	if r >= 3 && r < 6 {
		severity = "medium"
	} else if r > 6 {
		severity = "high"
	}

	sources := []string{"grafana", "manual"}
	regions := []string{"us-east", "us-west", "africa", "asia-pacific", "canada"}
	groups := []string{"ecommerce", "food", "map", "mail", "vpn"}
	teams := []string{"orderly", "paymate", "Invento", "reco", "engage"}
	products := []string{"orders", "warehouse", "alert", "rewards", "shipper"}

	startedAt := gutil.NewInt64(randDateInMilli(OneDayInMillis*10) - OneDayInMillis)
	seenAt := gutil.NewInt64(gutil.RandNumber(*startedAt, *startedAt+(8*60*1000)))
	endedAt := gutil.NewInt64(gutil.RandNumber(*seenAt, *seenAt+(25*60*1000)))
	return &CreateTicket{
		Fingerprint: gutil.RandString(20),
		Source:      gutil.NewString(sources[gutil.RandNumber(0, int64(len(sources)))]),
		Raw:         gutil.NewString(raw),
		Annotations: map[string]string{
			"name":         "high memory usage",
			"second value": "hi from seeder",
		},
		IsFiring:  gutil.NewBool(false),
		StartedAt: startedAt,
		EndedAt:   endedAt,
		Values: map[string]string{
			"a": "b",
		},
		GeneratorUrl: gutil.NewString("http://localhost:4000/fake_generator_url"),
		IsSpam:       gutil.NewBool(gutil.RandNumber(0, 10) >= 5),
		Severity:     gutil.NewString(severity),
		Title:        gutil.NewString("High memory usage from central region"),
		Description:  gutil.NewString("High memory usage from central region"),
		SeenAt:       seenAt,
		Labels: map[string]string{
			"label-a": "val-a",
			"region":  regions[gutil.RandNumber(0, int64(len(regions)))],
			"group":   groups[gutil.RandNumber(0, int64(len(groups)))],
			"team":    teams[gutil.RandNumber(0, int64(len(teams)))],
			"product": products[gutil.RandNumber(0, int64(len(products)))],
		},
	}
}

const OneDayInMillis = 24 * 60 * 60 * 1000

func randDateInMilli(randDiff int64) int64 {
	return time.Now().UnixMilli() - gutil.RandNumber(0, randDiff)
}
