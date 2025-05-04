package libvore

import "github.com/jmeaster30/vore/libvore/engine"

type (
	Matches     engine.Matches
	Match       engine.Match
	ReplaceMode engine.ReplaceMode
)

const (
	OVERWRITE ReplaceMode = ReplaceMode(engine.OVERWRITE)
	CONFIRM   ReplaceMode = ReplaceMode(engine.CONFIRM)
	NEW       ReplaceMode = ReplaceMode(engine.NEW)
	NOTHING   ReplaceMode = ReplaceMode(engine.NOTHING)
)
