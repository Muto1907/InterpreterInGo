package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hi := &String{Value: "Hi"}
	hi2 := &String{Value: "Hi"}
	other := &String{Value: "other"}
	other2 := &String{Value: "other"}

	if hi.HashKey() != hi2.HashKey() {
		t.Errorf("Different Hashkey on same-content-strings")
	}
	if other.HashKey() != other2.HashKey() {
		t.Errorf("Different Hashkey on same-content-strings")
	}
	if other.HashKey() == hi.HashKey() {
		t.Error("Same Hashkey on different-content-strings")
	}
}
