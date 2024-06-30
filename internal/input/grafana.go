package input

import (
	"fmt"
	"time"

	"github.com/kamva/gutil"
)

type GrafanaWebhookPayload struct {
	Alerts []GrafanaAlert `json:"alerts"`
}

type GrafanaAlert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	FingerPrint  string            `json:"fingerprint"`
	GeneratorURL string            `json:"generatorURL"`
	Values       map[string]any    `json:"values"`
}

// ToPatchTickets converts the input to PatchTicket list.
func (in *GrafanaWebhookPayload) ToPatchTickets() []*PatchTicket {
	res := make([]*PatchTicket, len(in.Alerts))
	for idx, val := range in.Alerts {
		res[idx] = val.ToPatchTicket()
	}
	return res
}

func (in *GrafanaAlert) ToPatchTicket() *PatchTicket {
	// Send the patch ticket input for firing state:
	values := make(map[string]string)
	for k, v := range in.Values {
		values[k] = fmt.Sprint(v)
	}

	// Mark grafana nodata and datasourceError alerts as spam.
	// TODO: convert status available values to constant.
	isSpam := gutil.NewBool(in.Labels["alertname"] == "DatasourceError" || in.Labels["alertname"] == "DatasourceNoData")

	var endedAt *int64
	if in.Status == "resolved" {
		endedAt = gutil.NewInt64(in.EndsAt.UnixMilli())
	}

	return &PatchTicket{
		IsFromTicketGenerator: true,
		Source:                gutil.NewString("grafana"),
		Fingerprint:           in.FingerPrint,
		Annotations:           in.Annotations,
		SyncAnnotations:       false,
		IsFiring:              gutil.NewBool(in.Status == "firing"),
		StartedAt:             gutil.NewInt64(in.StartsAt.UnixMilli()),
		Values:                values,
		GeneratorUrl:          &in.GeneratorURL,
		IsSpam:                isSpam,
		Title:                 gutil.NewString(gutil.StringDefault(in.Labels["alertname"], "(no_title)")),
		Labels:                in.Labels,
		SyncLabels:            false,
		EndedAt:               endedAt,
	}
}
