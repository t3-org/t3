package input

import (
	"fmt"
	"time"

	"github.com/kamva/gutil"
)

type GrafanaWebhookPayload struct {
	Alerts []GrafanaAlert
}

type GrafanaAlert struct {
	Status       string
	Labels       map[string]string
	Annotations  map[string]string
	StartsAt     time.Time
	EndsAt       time.Time
	FingerPrint  string
	generatorURL string
	values       map[string]any
}

func (in *GrafanaWebhookPayload) PatchInputs() []*PatchTicket {
	res := make([]*PatchTicket, len(in.Alerts))
	for idx, val := range in.Alerts {
		res[idx] = val.PatchInput()
	}
	return res
}

func (in *GrafanaAlert) PatchInput() *PatchTicket {
	if in.Status == "resolved" {
		return &PatchTicket{
			Fingerprint: in.FingerPrint,
			IsFiring:    gutil.NewBool(false),
			EndedAt:     gutil.NewInt64(in.EndsAt.UnixMilli()),
		}
	}

	// Send the patch ticket input for firing state:
	var values map[string]string
	for k, v := range in.values {
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
		GeneratorUrl:    &in.generatorURL,
		IsSpam:          gutil.NewBool(false),
		Labels:          in.Labels,
		SyncLabels:      false,
	}
}
