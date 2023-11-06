package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringMap map[string]string

func (a StringMap) Scan(value any) error {
	if value == nil {
		a = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &a)
	case string:
		return json.Unmarshal([]byte(v), &a)
	default:
		return fmt.Errorf("JSONText.Scan: expected []byte or string, got %T (%q)", value, value)
	}
}

func (a StringMap) Value() (driver.Value, error) {
	return json.Marshal(a)
}

var _ driver.Valuer = StringMap{}
var _ sql.Scanner = StringMap{}
