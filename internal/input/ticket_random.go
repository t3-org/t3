package input

import (
	"time"

	"github.com/kamva/gutil"
)

func RandomCreatTicket() *CreateTicket {
	severity := "low"
	r := gutil.RandNumber(0, 10)
	if r >= 3 && r < 6 {
		severity = "medium"
	} else if r > 6 {
		severity = "high"
	}

	sources := []string{"grafana", "manual"}
	regions := []string{"region-1", "region-2", "region-3", "region-4"}
	groups := []string{"group-a", "group-b", "group-c", "group-d", "group-e", "group-f"}
	teams := []string{"team-a", "team-b", "team-c", "team-d", "team-e", "team-f"}
	products := []string{"product-a", "product-b", "product-c", "product-d", "product-e", "product-f"}

	return &CreateTicket{
		Fingerprint: gutil.RandString(20),
		Source:      gutil.NewString(sources[gutil.RandNumber(0, int64(len(sources)))]),
		Raw:         nil,
		Annotations: map[string]string{
			"name":         "high memory usage",
			"second value": "hi from seeder",
		},
		IsFiring:  gutil.NewBool(false),
		StartedAt: gutil.NewInt64(randDateInMilli(OneDayInMillis) - OneDayInMillis),
		EndedAt:   gutil.NewInt64(randDateInMilli(OneDayInMillis) - OneDayInMillis/2),
		Values: map[string]string{
			"a": "b",
		},
		GeneratorUrl: gutil.NewString("http://localhost:4000/fake_generator_url"),
		IsSpam:       gutil.NewBool(gutil.RandNumber(0, 10) >= 5),
		Severity:     gutil.NewString(severity),
		Title:        gutil.NewString("High memory usage from central region"),
		Description:  gutil.NewString("High memory usage from central region"),
		SeenAt:       gutil.NewInt64(randDateInMilli(OneDayInMillis / 2)),
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
