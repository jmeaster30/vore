package engine

import "encoding/json"

type ValueType int

const (
	ValueStringType ValueType = iota
	ValueHashMapType
)

type Value interface {
	String() ValueString
	Hashmap() ValueHashMap
	Copy() Value
	ToGo() any
	MarshalJSON() ([]byte, error)

	getType() ValueType
	// add an interface that runs provided functions on each type of value
	process(matchPrintDepth int, hashmapFunc func(int, ValueHashMap), stringFunc func(ValueString))
}

type ValueString struct {
	Value string
}

func NewValueString(value string) ValueString {
	return ValueString{value}
}

func (v ValueString) ToGo() any {
	return v.Value
}

func (v ValueString) String() ValueString {
	return v
}

func (v ValueString) Hashmap() ValueHashMap {
	result := ValueHashMap{
		Value: make(map[string]Value),
	}

	result.Value["value"] = NewValueString(v.Value)
	return result
}

func (v ValueString) Copy() Value {
	return NewValueString(v.Value)
}

func (v ValueString) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

func (v ValueString) getType() ValueType {
	return ValueStringType
}

func (v ValueString) process(matchPrintDepth int, hashmapFunc func(int, ValueHashMap), stringFunc func(ValueString)) {
	stringFunc(v)
}

type ValueHashMap struct {
	Value map[string]Value
}

func NewValueHashMap() ValueHashMap {
	return ValueHashMap{
		Value: make(map[string]Value),
	}
}

func (v ValueHashMap) ToGo() any {
	result := make(map[string]any)
	for key, value := range v.Value {
		result[key] = value.ToGo()
	}
	return result
}

func (v ValueHashMap) String() ValueString {
	return NewValueString("[ValueHashMap]")
}

func (v ValueHashMap) Hashmap() ValueHashMap {
	return v
}

func (v ValueHashMap) Copy() Value {
	result := NewValueHashMap()
	for k, val := range v.Value {
		result.Add(k, val.Copy())
	}
	return result
}

func (v ValueHashMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

func (v ValueHashMap) getType() ValueType {
	return ValueHashMapType
}

func (v ValueHashMap) process(matchPrintDepth int, hashmapFunc func(int, ValueHashMap), stringFunc func(ValueString)) {
	hashmapFunc(matchPrintDepth, v)
}

func (v ValueHashMap) Get(name string) (Value, bool) {
	val, found := v.Value[name]
	return val, found
}

func (v ValueHashMap) Add(name string, value Value) {
	v.Value[name] = value
}

func (v ValueHashMap) Len() int {
	return len(v.Value)
}

func (v ValueHashMap) Keys() []string {
	res := []string{}
	for k := range v.Value {
		res = append(res, k)
	}
	return res
}

// TODO convert this to clean all characters of the new ValueHashMap
func cleanControlCharacters(s string) string {
	result := ""
	for _, c := range s {
		switch c {
		case '\t':
			result += "\\t"
		case '\r':
			result += "\\r"
		case '\n':
			result += "\\n"
		case '\\':
			result += "\\\\"
		case '\a':
			result += "\\a"
		case '\b':
			result += "\\b"
		case '\f':
			result += "\\f"
		case '\v':
			result += "\\v"
		default:
			result += string(c)
		}
	}
	return result
}
