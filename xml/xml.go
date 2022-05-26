package xml

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"util"
)

var (
	//ErrInvalidDocument 无效文档错误
	ErrXmlInvalidDocument = errors.New("invalid document")

	//ErrInvalidRoot 根级别的数据无效 err
	ErrXmlInvalidRoot = errors.New("data at the root level is invalid")
	//html 头
	XmlHeader string = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>`
	//map转xml 切片外层名称
	XmlChildName string = "item"
)

const (
	attrPrefix = "@"
	textPrefix = "#text"
)

// 解码器实例
type XmlDecoder struct {
	doc        string
	attrPrefix string
	textPrefix string
	recast     bool //如果 'recast' 是 'true'，那么如果可能，值将被转换为 boolean 或 float64。
}
type XmlNode struct {
	dup   bool   // is member of a list
	attr  bool   // is an attribute
	key   string // XML tag
	val   string // element value
	nodes []*XmlNode
}

// NewXmlDecoder 创建新的解码器实例
func NewXmlDecoder(reader string, recast ...bool) *XmlDecoder {
	var r bool
	if len(recast) > 0 {
		r = recast[0]
	}
	return NewXmlDecoderWithPrefix(reader, attrPrefix, textPrefix, r)
}

// NewXmlDecoderWithPrefix 使用自定义属性前缀和文本前缀创建新的解码器实例
func NewXmlDecoderWithPrefix(doc, attrPrefix, textPrefix string, r bool) *XmlDecoder {
	return &XmlDecoder{doc: doc, attrPrefix: attrPrefix, textPrefix: textPrefix, recast: r}
}
func (d *XmlDecoder) Unmarshal() (map[string]interface{}, error) {

	//将 XML 文档转换为节点树。
	//xml.Decoder 无法正确处理某些文档中的空格
	doc := regexp.MustCompile("[ \t\n\r]*<").ReplaceAllString(d.doc, "<")
	b := bytes.NewBufferString(doc)
	p := xml.NewDecoder(b)
	n, berr := xmlToTree("", nil, p, d.attrPrefix)
	if berr != nil {
		return nil, berr
	}
	m := make(map[string]interface{})
	m[n.key] = n.treeToMap(d.recast)

	return m, nil
}

// xmlToTree - xmlToTree - 将“干净”的 XML 文档加载到 *XmlNode 的树中。
func xmlToTree(skey string, a []xml.Attr, p *xml.Decoder, attrPrefix string) (*XmlNode, error) {
	n := new(XmlNode)
	n.nodes = make([]*XmlNode, 0)

	if skey != "" {
		n.key = skey
		if len(a) > 0 {
			for _, v := range a {
				na := new(XmlNode)
				na.attr = true
				na.key = `-` + v.Name.Local
				na.val = v.Value
				n.nodes = append(n.nodes, na)
			}
		}
	}
	for {
		t, err := p.Token()
		if err != nil {
			return nil, err
		}
		switch t := t.(type) {
		case xml.StartElement:
			tt := t
			// handle root
			if n.key == "" {
				n.key = tt.Name.Local
				if len(tt.Attr) > 0 {
					for _, v := range tt.Attr {
						na := new(XmlNode)
						na.attr = true
						na.key = attrPrefix + v.Name.Local
						na.val = v.Value
						n.nodes = append(n.nodes, na)
					}
				}
			} else {
				nn, nnerr := xmlToTree(tt.Name.Local, tt.Attr, p, attrPrefix)
				if nnerr != nil {
					return nil, nnerr
				}
				n.nodes = append(n.nodes, nn)
			}
		case xml.EndElement:
			// scan n.nodes for duplicate n.key values
			n.markDuplicateKeys()
			return n, nil
		case xml.CharData:
			tt := string(t)
			if len(n.nodes) > 0 {
				nn := new(XmlNode)
				nn.key = "#text"
				nn.val = tt
				n.nodes = append(n.nodes, nn)
			} else {
				n.val = tt
			}
		default:
			// noop
		}
	}
	// Logically we can't get here, but provide an error message anyway.
	return nil, errors.New("Unknown parse error in xmlToTree() for: " + n.key)
}

// (*XmlNode)markDuplicateKeys - 为加载 map[string]interface{} 设置 node.dup 标志。
func (n *XmlNode) markDuplicateKeys() {
	l := len(n.nodes)
	for i := 0; i < l; i++ {
		if n.nodes[i].dup {
			continue
		}
		for j := i + 1; j < l; j++ {
			if n.nodes[i].key == n.nodes[j].key {
				n.nodes[i].dup = true
				n.nodes[j].dup = true
			}
		}
	}
}

// (*XmlNode)treeToMap - 将节点树转换为 map[string]interface{}。
//（解析为与来自 json.Unmarshal() 的结构相同的映射。）
// 注意：root 没有实例化； 调用：“m[n.key] = n.treeToMap(recast)”。
func (n *XmlNode) treeToMap(r bool) interface{} {
	if len(n.nodes) == 0 {
		return util.Recast(n.val, r)
	}

	m := make(map[string]interface{}, 0)
	for _, v := range n.nodes {
		// just a value
		if !v.dup && len(v.nodes) == 0 {
			m[v.key] = util.Recast(v.val, r)
			continue
		}

		// a list of values
		if v.dup {
			var a []interface{}
			if vv, ok := m[v.key]; ok {
				a = vv.([]interface{})
			} else {
				a = make([]interface{}, 0)
			}
			a = append(a, v.treeToMap(r))
			m[v.key] = interface{}(a)
			continue
		}

		// it's a unique key
		m[v.key] = v.treeToMap(r)
	}

	return interface{}(m)
}

// WriteMap - 转储 map[string]interface{} 以供检查。
// 'offset' 是初始缩进计数； 通常：WriteMap(m)。
// 注意：对于 XML，所有元素类型都是“字符串”。
// 但是代码编写为通用代码，用于 json.Unmarshal() 中的 maps[string]interface{} 值。
func WriteMap(m interface{}, offset ...int) string {
	var indent int
	if len(offset) == 1 {
		indent = offset[0]
	}

	var s string
	switch m := m.(type) {
	case nil:
		return "[nil] nil"
	case string:
		return "[string] " + m
	case float64:
		return "[float64] " + strconv.FormatFloat(m, 'e', 2, 64)
	case bool:
		return "[bool] " + strconv.FormatBool(m)
	case []interface{}:
		s += "[[]interface{}]"
		for i, v := range m {
			s += "\n"
			for i := 0; i < indent; i++ {
				s += "  "
			}
			s += "[item: " + strconv.FormatInt(int64(i), 10) + "]"
			switch v.(type) {
			case string, float64, bool:
				s += "\n"
			default:
				// noop
			}
			for i := 0; i < indent; i++ {
				s += "  "
			}
			s += WriteMap(v, indent+1)
		}
	case map[string]interface{}:
		for k, v := range m {
			s += "\n"
			for i := 0; i < indent; i++ {
				s += "  "
			}
			// s += "[map[string]interface{}] "+k+" :"+WriteMap(v,indent+1)
			s += k + " :" + WriteMap(v, indent+1)
		}
	default:
		// shouldn't ever be here ...
		s += fmt.Sprintf("unknown type for: %v", m)
	}
	return s
}

type XmlIndentation struct {
	Prefix string
	Indent string
}

type XmlRoot struct {
	Name          *xml.Name
	XMLAttributes *[]xml.Attr
	Attributes    map[string]string
}
type XmlStructMap struct {
	CData   bool
	Map     map[string]interface{}
	Indent  *XmlIndentation
	XmlRoot *XmlRoot
}

type xmlMapEntry struct {
	XMLName    xml.Name
	Value      interface{} `xml:",innerxml"`
	CDataValue interface{} `xml:",cdata"`
}

//NewXmlEncoder 初始化map转xml实例
func NewXmlEncoder(input map[string]interface{}) *XmlStructMap {
	return &XmlStructMap{Map: input}
}

//WithIndent xml 添加缩进 WithIndent("", "  ")
func (smap *XmlStructMap) WithIndent(prefix string, indent string) *XmlStructMap {
	smap.Indent = &XmlIndentation{Prefix: prefix, Indent: indent}
	return smap
}

//将根节点添加到 xml  WithXmlRoot("person", map[string]string{"mood": "happy"})
func (smap *XmlStructMap) WithXmlRoot(name string, attributes map[string]string) *XmlStructMap {
	attr := []xml.Attr{}
	for k, v := range attributes {
		attr = append(attr, xml.Attr{Name: xml.Name{Local: k}, Value: v})
	}
	smap.XmlRoot = &XmlRoot{Name: &xml.Name{Local: name}, XMLAttributes: &attr, Attributes: attributes}
	return smap
}

//CDATA 标记添加到所有节点
func (smap *XmlStructMap) AsCData() *XmlStructMap {
	smap.CData = true
	return smap
}

//Print 以 json 格式打印配置
func (smap *XmlStructMap) Print() *XmlStructMap {
	var indent interface{} = smap.Indent
	var root interface{}
	if smap.Indent != nil {
		indent = map[string]int{"indent_spaces": len(smap.Indent.Indent), "prefix_spaces": len(smap.Indent.Prefix)}
	}
	// if root = smap.XmlRoot; root != nil {
	root = map[string]interface{}{"name": smap.XmlRoot.Name.Local, "attributes": smap.XmlRoot.Attributes}
	// }
	b, _ := json.MarshalIndent(map[string]interface{}{"root": root, "cdata": smap.CData, "indent": indent}, " ", "  ")
	fmt.Println(string(b))
	return smap
}

//Builds XML as bytes
func (smap *XmlStructMap) Marshal() ([]byte, error) {
	var (
		xmlBytes []byte
		err      error
	)
	if smap.Indent == nil {
		xmlBytes, err = xml.Marshal(smap)
	} else {
		xmlBytes, err = xml.MarshalIndent(smap, smap.Indent.Prefix, smap.Indent.Indent)
	}
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer([]byte(XmlHeader + "\n"))
	buff.Write(xmlBytes)
	return buff.Bytes(), nil
}

//Builds XML as string
func (smap *XmlStructMap) MarshalToString() (string, error) {
	xmlBytes, err := smap.Marshal()
	return string(xmlBytes), err
}

func (smap XmlStructMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	if len(smap.Map) == 0 {
		return nil
	}
	if smap.XmlRoot != nil {
		start = xml.StartElement{Name: *smap.XmlRoot.Name, Attr: *smap.XmlRoot.XMLAttributes}
		if err := e.EncodeToken(start); err != nil {
			return err
		}
	}

	for k, v := range smap.Map {
		if err := handleXmlChildren(e, k, v, smap.CData); err != nil {
			return err
		}
	}

	if smap.XmlRoot != nil {
		return e.EncodeToken(start.End())
	}
	return nil
}

func handleXmlChildren(e *xml.Encoder, fieldName string, v interface{}, cdata bool) error {
	//返回数字map 补充后的keys 和是否是连续intmap
	var getKeyConseIntMapAndAdd func(map[string]interface{}) ([]string, bool) = func(m map[string]interface{}) ([]string, bool) {
		var keys []string
		if len(m) == 0 {
			return keys, false
		}
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys) //升序

		var newKeys []string
		var dk int64 = -1 //起始下标
		for _, k := range keys {
			gk, err := strconv.ParseInt(k, 10, 32)
			if err != nil {
				return keys, false
			}
			if dk = dk + 1; dk != gk { //不是连续的int 补全key
				for i := dk; i <= int64(math.Abs(float64(gk-dk))); i++ {
					newKeys = append(newKeys, strconv.Itoa(int(i)))
				}
				dk = gk
			} else {
				newKeys = append(newKeys, k)
			}
		}
		return newKeys, true
	}

	if reflect.TypeOf(v) == nil {
		return e.Encode(xmlMapEntry{XMLName: xml.Name{Local: fieldName}, Value: ""})
	} else if reflect.TypeOf(v).Kind() == reflect.Map {
		e.EncodeToken(xml.StartElement{Name: xml.Name{Local: fieldName}})
		switch v := v.(type) {
		case map[string]interface{}:

			//填充key 解析
			keys, isKeyConseInt := getKeyConseIntMapAndAdd(v)
			for _, key := range keys {
				if key == "xml_child_name" {
					continue
				}
				childKey := key
				if isKeyConseInt {
					childKey = XmlChildName
				}
				val, ok := v[key]
				if !ok {
					val = ""
				}
				// fmt.Println("old", reflect.TypeOf(val), val)
				if reflect.TypeOf(val).Kind() == reflect.Map {
					//EncodeToken  一个外层标签
					e.EncodeToken(xml.StartElement{Name: xml.Name{Local: childKey}})

					for mk, vd := range val.(map[string]interface{}) {
						kind := reflect.TypeOf(vd).Kind()
						if kind == reflect.Slice {
							for _, elem := range vd.([]interface{}) {
								handleXmlChildren(e, mk, elem, cdata)
							}
						} else {
							handleXmlChildren(e, mk, vd, cdata)
						}
					}
					e.EncodeToken(xml.EndElement{Name: xml.Name{Local: childKey}})

				} else {
					handleXmlChildren(e, childKey, val, cdata)
				}

				//不需要填充key 直接解析
				// for key, val := range v {
				// 	if key == "xml_child_name" {
				// 		continue
				// 	}
				// 	handleXmlChildren(e, key, val, cdata)
				// }
			}

		case []string:
			for _, val := range v {
				handleXmlChildren(e, fieldName, val, cdata)
			}
		case []interface{}:
			for _, val := range v {
				handleXmlChildren(e, fieldName, val, cdata)
			}
		}

		return e.EncodeToken(xml.EndElement{Name: xml.Name{Local: fieldName}})
	} else if reflect.TypeOf(v).Kind() == reflect.Slice {
		e.EncodeToken(xml.StartElement{Name: xml.Name{Local: fieldName}})
		switch v := v.(type) {
		case []map[string]interface{}:
			if _, hasChildName := v[0]["xml_child_name"]; hasChildName {
				XmlChildName = v[0]["xml_child_name"].(string)
			}
			for _, elem := range v {
				for _, pv := range elem {
					handleXmlChildren(e, XmlChildName, pv, cdata)
				}
			}
		case []string:
			for _, elem := range v {
				handleXmlChildren(e, XmlChildName, elem, cdata)
			}
		case []interface{}:
			for _, elem := range v {
				handleXmlChildren(e, XmlChildName, elem, cdata)
			}
		default:
		}

		return e.EncodeToken(xml.EndElement{Name: xml.Name{Local: fieldName}})
	}
	if cdata {
		return e.Encode(xmlMapEntry{XMLName: xml.Name{Local: fieldName}, CDataValue: v})
	} else {
		return e.Encode(xmlMapEntry{XMLName: xml.Name{Local: fieldName}, Value: fmt.Sprintf("%v", v)})
	}
}
