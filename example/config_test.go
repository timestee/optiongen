package example

import "testing"

// interface test
type interfaceTest struct {
	data []int64
}

func (i *interfaceTest) Apply(cc *Config) ConfigOption {
	return AppendTestSliceInt64(1, 2, 3, 4).Apply(cc)
}

func TestNewConfig(t *testing.T) {
	it := &interfaceTest{data: []int64{1, 2, 3, 4}}
	tc := NewFuncNameSpecified(false, "", WithTestMapIntInt(map[int]int{2: 4}), it)
	if tc == nil {
		t.Fatal("new config error")
	}
	tc.GetFOO()
	if tc.GetTestMapIntInt()[2] != 4 {
		t.Fatal("map get val error")
	}
	previousValue := tc.GetTestInt()
	changeTo := 1232323232323232
	previous := tc.ApplyOption(WithTestInt(changeTo))
	if tc.GetTestInt() != changeTo {
		t.Fatal("ApplyOption failed")
	}
	tc.ApplyOption(previous...)
	if tc.GetTestInt() != previousValue {
		t.Fatal("ApplyOption Restore failed")
	}

	WithXXXXXXRedis(nil)
}
