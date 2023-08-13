package pkpu

import "testing"

func TestGetDataPkpu(t *testing.T) {
	x := 10
	res := GetDataPkpu(x)
	if res == nil {
		t.Error()
		return
	}
}
