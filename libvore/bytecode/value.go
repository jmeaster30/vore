package bytecode

import "strconv"

type ValueType int

const (
	ValueType_String ValueType = iota
	ValueType_Number
	ValueType_Boolean
)

func (vt ValueType) String() string {
	switch vt {
	case ValueType_String:
		return "str"
	case ValueType_Number:
		return "num"
	case ValueType_Boolean:
		return "bool"
	}
	return "unk"
}

type Value interface {
	GetType() ValueType
	GetString() string
	GetNumber() int
	GetBoolean() bool
}

type StringValue struct {
	value string
}

func NewString(value string) Value {
	return StringValue{value}
}

func (v StringValue) GetString() string {
	return v.value
}

func (v StringValue) GetNumber() int {
	intval, err := strconv.Atoi(v.value)
	if err != nil {
		return 0
	}
	return intval
}

func (v StringValue) GetBoolean() bool {
	return len(v.value) != 0
}

func (v StringValue) GetType() ValueType {
	return ValueType_String
}

type NumberValue struct {
	value int
}

func NewNumber(value int) NumberValue {
	return NumberValue{value}
}

func (v NumberValue) GetString() string {
	return strconv.Itoa(v.value)
}

func (v NumberValue) GetNumber() int {
	return v.value
}

func (v NumberValue) GetBoolean() bool {
	return v.value != 0
}

func (v NumberValue) GetType() ValueType {
	return ValueType_Number
}

type BooleanValue struct {
	value bool
}

func NewBoolean(value bool) BooleanValue {
	return BooleanValue{value}
}

func (v BooleanValue) GetString() string {
	if v.value {
		return "true"
	}
	return "false"
}

func (v BooleanValue) GetNumber() int {
	if v.value {
		return 1
	}
	return 0
}

func (v BooleanValue) GetBoolean() bool {
	return v.value
}

func (v BooleanValue) GetType() ValueType {
	return ValueType_Boolean
}
