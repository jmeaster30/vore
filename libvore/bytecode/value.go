package bytecode

import "strconv"

type ValueType int

const (
	ValueType_String ValueType = iota
	ValueType_Number
	ValueType_Boolean
)

type Value interface {
	GetType() ValueType
	GetString() string
	GetNumber() int
	GetBoolean() bool
}

type StringValue struct {
	Value string
}

func (v StringValue) GetString() string {
	return v.Value
}

func (v StringValue) GetNumber() int {
	intval, err := strconv.Atoi(v.Value)
	if err != nil {
		return 0
	}
	return intval
}

func (v StringValue) GetBoolean() bool {
	return len(v.Value) != 0
}

func (v StringValue) GetType() ValueType {
	return ValueType_String
}

type NumberValue struct {
	Value int
}

func (v NumberValue) GetString() string {
	return strconv.Itoa(v.Value)
}

func (v NumberValue) GetNumber() int {
	return v.Value
}

func (v NumberValue) GetBoolean() bool {
	return v.Value != 0
}

func (v NumberValue) GetType() ValueType {
	return ValueType_Number
}

type BooleanValue struct {
	Value bool
}

func (v BooleanValue) GetString() string {
	if v.Value {
		return "true"
	}
	return "false"
}

func (v BooleanValue) GetNumber() int {
	if v.Value {
		return 1
	}
	return 0
}

func (v BooleanValue) GetBoolean() bool {
	return v.Value
}

func (v BooleanValue) GetType() ValueType {
	return ValueType_Boolean
}
