package engine

import (
	"encoding/json"
	"fmt"

	"github.com/jmeaster30/vore/libvore/bytecode"
	"github.com/jmeaster30/vore/libvore/ds"
)

type Matches []Match

func (m Matches) Print() {
	fmt.Println("============")
	for _, match := range m {
		match.Print()
		fmt.Println("============")
	}
}

func (m Matches) Json() string {
	var mi any = m
	data, err := json.Marshal(mi.([]Match))
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (m Matches) FormattedJson() string {
	data, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(data)
}

type Match struct {
	Filename    string
	MatchNumber int
	Offset      ds.Range
	Line        ds.Range
	Column      ds.Range
	Value       string
	Replacement ds.Optional[string]
	Variables   bytecode.MapValue
}

func (m Match) Json() string {
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (m Match) FormattedJson() string {
	data, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (m Match) MarshalJSON() ([]byte, error) {
	result := make(map[string]any)
	result["filename"] = m.Filename
	result["matchNumber"] = m.MatchNumber
	result["offset"] = m.Offset
	result["line"] = m.Line
	result["column"] = m.Column
	result["value"] = m.Value
	if m.Replacement.HasValue() {
		result["replacement"] = m.Replacement.GetValue()
	}
	result["variables"] = m.Variables
	return json.Marshal(result)
}

func (m Match) Print() {
	fmt.Printf("Filename: %s\n", m.Filename)
	fmt.Printf("MatchNumber: %d\n", m.MatchNumber)
	fmt.Printf("Value: '%s'\n", m.Value)
	if m.Replacement.HasValue() {
		fmt.Printf("Replaced: %s\n", m.Replacement.GetValue())
	}
	fmt.Printf("Offset: %d %d\n", m.Offset.Start, m.Offset.End)
	fmt.Printf("Line: %d %d\n", m.Line.Start, m.Line.End)
	fmt.Printf("Column: %d %d\n", m.Column.Start, m.Column.End)
	fmt.Println("Variables:")
	fmt.Println("  [key] = [value]")

	for key, value := range m.Variables.Map() {
		fmt.Printf("  %s = %s\n", key, value)
	}

	fmt.Println()
}
