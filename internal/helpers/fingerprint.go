package helpers

import (
	"fmt"
)

// GlobalFingerprint generates a global fingerprint using
// started_at(preferably send in milliseconds) and fingerprint.
func GlobalFingerprint(startedAt *int64, fingerprint string) string {
	if startedAt == nil || fingerprint == "" {
		return ""
	}
	return fmt.Sprintf("%d-%s", *startedAt, fingerprint)
}
