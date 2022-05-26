package util

import (
	"encoding/json"
	"testing"
)

func TestParseStr(t *testing.T) {
	result, err := ParseStr(`a=val&a[0][3]=2.5`)

	if err != nil {
		t.Log(err)
	}

	sb, err := json.Marshal(result)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(sb))

}
