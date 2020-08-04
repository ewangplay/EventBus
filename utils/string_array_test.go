package utils

import (
	"testing"
)

func TestStringArray(t *testing.T) {
	var arr StringArray
	arr.Set("hello")
	arr.Set("world")
	if len(arr) != 2 {
		t.FailNow()
	}
	t.Logf("%v", arr)
}
