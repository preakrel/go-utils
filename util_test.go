package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseStr(t *testing.T) {
	result, err := ParseStr(`poc=true&list[0]=1&list[1]=2&list[2][0]=3&list[2][1]=4&list[2][2][0]=5&list[2][2][1]=6&list[2][2][2][0]=7&list[2][2][2][1]=8&list[2][2][2][2][0]=9&list[2][2][2][2][1]=10&list[2][2][2][2][2][0]=11&list[2][2][2][2][2][1]=12&list[2][2][2][2][2][2][0]=13&list[2][2][2][2][2][2][1]=14&id=187923`)

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

	jsonStr := `{"id":187923,"poc":true,"list":[1,2,[3,4,[5,6,[7,8,[9,10,[11,12,[13,14]]]]]]],"map":{"m1":"k1","m2":"k2"}}`
	var r map[string]interface{}
	d := json.NewDecoder(bytes.NewReader([]byte(jsonStr)))
	d.UseNumber()
	err := d.Decode(&r)

	if err != nil {
		t.Log(err)
	}
	fmt.Println(r)
	str := HTTPBuildQuery(r, "", "&", "QUERY_RFC3986")

	t.Log(str)

}
