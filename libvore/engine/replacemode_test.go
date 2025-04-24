package engine

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestOverwriteString(t *testing.T) {
	mode := OVERWRITE
	testutils.AssertEqual(t, "OVERWRITE", mode.String())
}

func TestConfirmString(t *testing.T) {
	mode := CONFIRM
	testutils.AssertEqual(t, "CONFIRM", mode.String())
}

func TestNewString(t *testing.T) {
	mode := NEW
	testutils.AssertEqual(t, "NEW", mode.String())
}

func TestNothingString(t *testing.T) {
	mode := NOTHING
	testutils.AssertEqual(t, "NOTHING", mode.String())
}

func TestUnknownString(t *testing.T) {
	var mode ReplaceMode = 7
	testutils.AssertEqual(t, "", mode.String())
}
