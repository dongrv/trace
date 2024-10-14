package trace

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestTrace(t *testing.T) {
	t.Run("Pop", func(t *testing.T) {
		trace := Pop("sql").
			Set("SELECT * FROM users WHERE UID=%d", []interface{}{1}).
			SetDescription("err", errors.New("connection fail").Error()).
			Stop()
		defer clear(trace)
		if bytes, err := json.Marshal(trace); err != nil {
			t.Fail()
			return
		} else {
			t.Logf("%s", bytes)
		}
	})
}
