package cmds

import (
	"testing"
)

func TestUniqueOrderedSet(t *testing.T) {
	uos := UniqueOrderedSet{}

	uos = uos.Add("first")
	if uos[0] != "first" || len(uos) != 1 {
		t.Error("Adding single element to UniqueOrderedSet failed")
	}

	uos = uos.Add("first")
	if uos[0] != "first" || len(uos) != 1 {
		t.Error("Adding value to UniqueOrderedSet allowed duplicate")
	}

	uos = uos.Add("second")
	if uos[0] != "first" || uos[1] != "second" || len(uos) != 2 {
		t.Error("Adding second non duplicate value to UniqueOrderedSet failed")
	}
}
