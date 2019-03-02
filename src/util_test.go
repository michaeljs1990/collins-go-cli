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

func TestBytesToHumanSize(t *testing.T) {
	out := BytesToHumanSize(12)
	if out != "12 B" {
		t.Error("Bytes to human was suppose to be 12 B but was " + out)
	}

	out = BytesToHumanSize(1024)
	if out != "1 KB" {
		t.Error("Bytes to human was suppose to be 1 KB but was " + out)
	}

	out = BytesToHumanSize(1230)
	if out != "1.2 KB" {
		t.Error("Bytes to human was suppose to be 1.2 KB but was " + out)
	}

	out = BytesToHumanSize(146163105792)
	if out != "136.13 GB" {
		t.Error("Bytes to human was suppose to be 136.13 GB but was " + out)
	}

	out = BytesToHumanSize(1461631057920)
	if out != "1.33 TB" {
		t.Error("Bytes to human was suppose to be 1.33 TB but was " + out)
	}

	out = BytesToHumanSize(1461631057920000)
	if out != "1.3 PB" {
		t.Error("Bytes to human was suppose to be 1.3 PB but was " + out)
	}
}
