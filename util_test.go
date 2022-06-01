package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseStr(t *testing.T) {
	result, err := ParseStr(`id=187923&po1c[0]=1&po1c[1][0]=3&po1c[1][1][0]=5&po1c[1][1][1]=6&po1c[1][2]=4&po1c[2]=2&poc=true&test=abc`)

	if err != nil {
		t.Log(err)
	}

	sb, err := json.Marshal(result)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(sb))

}
func TestOriginalHttpBuildQuery(t *testing.T) {

	jsonStr := `{"id":187923,"po1c":[1,[3,[5,6],4],2],"poc":true,"test[0][1]":"2","test[0][3]":"12"}`
	var r map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &r)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(r)
	str := HTTPBuildQuery(r, "", "", "")

	t.Log(str)

}
