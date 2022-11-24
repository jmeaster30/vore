package libvore

type Range struct {
	Start int
	End   int
}

func NewRange(start int, end int) *Range {
	return &Range{start, end}
}
