package timer

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID        string          `json:"id"`
	ExecuteAt time.Time       `json:"execute_at"`
	Method    string          `json:"method"`
	URL       string          `json:"callback"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}
