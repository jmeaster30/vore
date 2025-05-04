package bytecode

import (
	"encoding/json"
	"strconv"

	"github.com/jmeaster30/vore/libvore/ds"
)

type ValueType int

const (
	ValueType_String ValueType = iota
	ValueType_Number
	ValueType_Boolean
	ValueType_Map
)

func (vt ValueType) String() string {
	switch vt {
	case ValueType_String:
		return "str"
	case ValueType_Number:
		return "num"
	case ValueType_Boolean:
		return "bool"
	case ValueType_Map:
		return "map"
	}
	panic("unknown value type")
}

type Value interface {
	Type() ValueType
	Any() any
	String() string
	Number() int
	Boolean() bool
	Map() map[string]any
	Copy() Value
}

func NewValue(value any) Value {
	switch v := value.(type) {
	case string:
		return NewString(v)
	case int:
		return NewNumber(v)
	case bool:
		return NewBoolean(v)
	case map[string]any:
		return NewMap(v)
	}
	panic("unknown value type")
}

func ToStringValue(value Value) MapValue {
	if value.Type() == ValueType_String {
		return value.(MapValue)
	}
	panic("cannot convert value to string value")
}

func ToNumberValue(value Value) NumberValue {
	if value.Type() == ValueType_Number {
		return value.(NumberValue)
	}
	panic("cannot convert value to number value")
}

func ToBooleanValue(value Value) BooleanValue {
	if value.Type() == ValueType_Boolean {
		return value.(BooleanValue)
	}
	panic("cannot convert value to boolean value")
}

func ToMapValue(value Value) MapValue {
	if value.Type() == ValueType_Boolean {
		return value.(MapValue)
	}
	panic("cannot convert value to map value")
}

type StringValue struct {
	value string
}

func NewString(value string) StringValue {
	return StringValue{value}
}

func (v StringValue) String() string {
	return v.value
}

func (v StringValue) Number() int {
	intval, err := strconv.Atoi(v.value)
	if err != nil {
		return 0
	}
	return intval
}

func (v StringValue) Boolean() bool {
	return len(v.value) != 0
}

func (v StringValue) Map() map[string]any {
	return map[string]any{"value": v.value}
}

func (v StringValue) Any() any {
	return v.value
}

func (v StringValue) Type() ValueType {
	return ValueType_String
}

func (v StringValue) Copy() Value {
	return NewString(v.value)
}

type NumberValue struct {
	value int
}

func NewNumber(value int) NumberValue {
	return NumberValue{value}
}

func (v NumberValue) String() string {
	return strconv.Itoa(v.value)
}

func (v NumberValue) Number() int {
	return v.value
}

func (v NumberValue) Boolean() bool {
	return v.value != 0
}

func (v NumberValue) Map() map[string]any {
	return map[string]any{"value": v.value}
}

func (v NumberValue) Any() any {
	return v.value
}

func (v NumberValue) Type() ValueType {
	return ValueType_Number
}

func (v NumberValue) Copy() Value {
	return NewNumber(v.value)
}

type BooleanValue struct {
	value bool
}

func NewBoolean(value bool) BooleanValue {
	return BooleanValue{value}
}

func (v BooleanValue) String() string {
	if v.value {
		return "true"
	}
	return "false"
}

func (v BooleanValue) Number() int {
	if v.value {
		return 1
	}
	return 0
}

func (v BooleanValue) Boolean() bool {
	return v.value
}

func (v BooleanValue) Map() map[string]any {
	return map[string]any{"value": v.value}
}

func (v BooleanValue) Any() any {
	return v.value
}

func (v BooleanValue) Type() ValueType {
	return ValueType_Boolean
}

func (v BooleanValue) Copy() Value {
	return NewBoolean(v.value)
}

type MapValue struct {
	value map[string]Value
}

func NewMap(value map[string]any) MapValue {
	result := make(map[string]Value)
	for key, val := range value {
		result[key] = NewValue(val)
	}
	return MapValue{result}
}

func NewEmptyMap() MapValue {
	return MapValue{make(map[string]Value)}
}

func (v MapValue) Get(index string) (Value, bool) {
	val, prs := v.value[index]
	return val, prs
}

func (v *MapValue) Set(index string, value Value) {
	v.value[index] = value
}

func (v *MapValue) Entries() []ds.Pair[string, Value] {
	result := make([]ds.Pair[string, Value], len(v.value))
	for key, value := range v.value {
		result = append(result, ds.NewPair(key, value))
	}
	return result
}

func (v MapValue) String() string {
	bytes, err := json.Marshal(v.Any())
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (v MapValue) Number() int {
	return 0
}

func (v MapValue) Boolean() bool {
	return len(v.value) != 0
}

func (v MapValue) Map() map[string]any {
	result := make(map[string]any)
	for key, value := range v.value {
		result[key] = value.Any()
	}
	return result
}

func (v MapValue) Any() any {
	result := make(map[string]any)
	for key, value := range v.value {
		result[key] = value.Any()
	}
	return result
}

func (v MapValue) Type() ValueType {
	return ValueType_Map
}

func (v MapValue) Copy() Value {
	return NewMap(v.Map())
}
