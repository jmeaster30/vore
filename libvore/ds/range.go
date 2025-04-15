package ds

import "encoding/json"

type Range struct {
	Start int
	End   int
}

func (r Range) MarshalJSON() ([]byte, error) {
	result := make(map[string]int)
	result["start"] = r.Start
	result["end"] = r.End
	return json.Marshal(result)
}

func NewRange(start int, end int) *Range {
	return &Range{start, end}
}
