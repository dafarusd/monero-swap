package swap2

import "testing"

func TestParseDecimalExact(t *testing.T) {
	cases := map[string]struct {
		dec  uint
		want string
	}{
		"0.004": {12, "4000000000"}, "1": {18, "1000000000000000000"}, "0.1": {6, "100000"}, "15.5": {12, "15500000000000"},
		"0.000000000001": {12, "1"}, "123456.789": {6, "123456789000"},
	}
	for in, c := range cases {
		got, err := ParseDecimal(in, c.dec)
		if err != nil || got.String() != c.want {
			t.Errorf("ParseDecimal(%q,%d) = %v, %v; want %s", in, c.dec, got, err, c.want)
		}
	}
	if _, err := ParseDecimal("0.0000001", 6); err == nil {
		t.Error("too many decimals should fail")
	}
	if _, err := ParseDecimal("-1", 6); err == nil {
		t.Error("negative should fail")
	}
	a := Asset{Symbol: "MUSD", Decimals: 6}
	if a.XmrPerAsset(4000000000).String() != "4000000000000000000000" {
		t.Errorf("XmrPerAsset: %s", a.XmrPerAsset(4000000000))
	}
}
