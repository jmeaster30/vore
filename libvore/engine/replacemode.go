package engine

type ReplaceMode int

const (
	OVERWRITE ReplaceMode = iota
	CONFIRM
	NEW
	NOTHING
)

func (r ReplaceMode) String() string {
	switch r {
	case OVERWRITE:
		return "OVERWRITE"
	case CONFIRM:
		return "CONFIRM"
	case NEW:
		return "NEW"
	case NOTHING:
		return "NOTHING"
	}
	return ""
}
