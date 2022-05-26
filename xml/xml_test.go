package xml

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestXml(t *testing.T) {

	data := `
<xml>
  <id>187923</id>
  <poc>
    <item>a</item>
    <item>b</item>
    <item>c</item>
    <item>
      <item>d</item>
    </item>
  </poc>
</xml>`
	m, err := NewXmlDecoder((data)).Unmarshal()
	if err != nil {
		t.Log(err)
	}

	jsonbytes, _ := json.Marshal(m)
	fmt.Println(string(jsonbytes))
	//map[xml:map[id:187923 poc:map[item:[c map[item:map[item:map[item:f]]]] item1:b item2:b item3:b]]]
	// str := OriginalHttpBuildQuery(m["xml"].(map[string]interface{}))
	// fmt.Println(m["xml"], str)
	// result := make(map[string]interface{})
	// StrParse("id=1&name=john&sub[]=1&sub[]=2&sub[]=three", result)
	// fmt.Println(result)
	// nResult := map[string]interface{}{
	// 	"xml": result,
	// }
	// var pm map[string]interface{}
	// err = json.Unmarshal([]byte(`{"a":{"1":["value"]}}`), &pm)
	// if err != nil {
	// 	t.Log(err)
	// }
	// fmt.Println(pm)
	str, err := NewXmlEncoder(m).WithIndent("", "  ").MarshalToString()
	if err != nil {
		t.Log(err)
	}
	// fmt.Println(pm)
	fmt.Println(str)
}
