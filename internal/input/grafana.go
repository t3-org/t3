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

func (in *GrafanaWebhookPayload) PatchInputs(ch *Channel) []*PatchTicket {
	res := make([]*PatchTicket, len(in.Alerts))
	for idx, val := range in.Alerts {
		res[idx] = val.PatchInput(ch)
	}
	return res
}

func (in *GrafanaAlert) PatchInput(ch *Channel) *PatchTicket {
	if in.Status == "resolved" {
		return &PatchTicket{
			Fingerprint: in.FingerPrint,
			IsFiring:    gutil.NewBool(false),
			EndedAt:     gutil.NewInt64(in.EndsAt.UnixMilli()),
			Channel:     ch,
		}
	}

	// Send the patch ticket input for firing state:
	var values map[string]string
	for k, v := range in.Values {
		values[k] = fmt.Sprint(v)
	}

	return &PatchTicket{
		Source:          gutil.NewString("grafana"),
		Fingerprint:     in.FingerPrint,
		Annotations:     in.Annotations,
		SyncAnnotations: false,
		IsFiring:        gutil.NewBool(in.Status == "firing"),
		StartedAt:       gutil.NewInt64(in.StartsAt.UnixMilli()),
		Values:          values,
		GeneratorUrl:    &in.GeneratorURL,
		IsSpam:          gutil.NewBool(false),
		Labels:          in.Labels,
		SyncLabels:      false,
		Channel:         ch,
	}
}