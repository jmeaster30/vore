package libvore

type Range struct {
	Start uint64
	End   uint64
}

func NewRange(start uint64, end uint64) *Range {
	return &Range{start, end}
}
