package pray_schedule

import "testing"

func TestGetDataPkpu(t *testing.T) {
	x := 10
	res := GetDataPkpu(x)
	if res == nil {
		t.Error()
		return
	}
}
