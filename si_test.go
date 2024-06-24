package si

import "testing"

func TestFormatAppend(t *testing.T) {
	var tests = []struct {
		V    int64
		Base Prefix
		Prec int
		Want string
	}{
		// Augment from Prefixed to No prefix.
		0: {V: 1000, Base: PrefixMilli, Prec: 4, Want: "1"},
		1: {V: 10000, Base: PrefixMilli, Prec: 4, Want: "10"},
		2: {V: 100000, Base: PrefixMilli, Prec: 4, Want: "100"},
		// Print within prefix range
		3: {V: 1, Base: PrefixMilli, Prec: 4, Want: "1m"},
		4: {V: 10, Base: PrefixMilli, Prec: 4, Want: "10m"},
		5: {V: 100, Base: PrefixMilli, Prec: 4, Want: "100m"},
		// Augment from low Prefixed to higher prefix.
		6: {V: 1_000_000, Base: PrefixMilli, Prec: 4, Want: "1k"},
		7: {V: 10_000_000, Base: PrefixMilli, Prec: 4, Want: "10k"},
		8: {V: 100_000_000, Base: PrefixMilli, Prec: 4, Want: "100k"},
		// Augment Prefix to No Prefix with decimals.
		9:  {V: 1234, Base: PrefixMilli, Prec: 4, Want: "1.234"},
		10: {V: 12344, Base: PrefixMilli, Prec: 4, Want: "12.34"},
		11: {V: 123444, Base: PrefixMilli, Prec: 4, Want: "123.4"},
		// Augment low to high prefix with decimals.
		12: {V: 1_234_444, Base: PrefixMilli, Prec: 4, Want: "1.234k"},
		13: {V: 12_344_444, Base: PrefixMilli, Prec: 4, Want: "12.34k"},
		14: {V: 123_444_444, Base: PrefixMilli, Prec: 4, Want: "123.4k"},
		// Augment low to high prefix with decimal chop-off.
		15: {V: 1_234_567, Base: PrefixMilli, Prec: 1, Want: "1k"},
		16: {V: 12_345_678, Base: PrefixMilli, Prec: 2, Want: "12k"},
		17: {V: 123_456_789, Base: PrefixMilli, Prec: 3, Want: "123k"},
		// Rounding simple.
		18: {V: 1500, Base: PrefixMilli, Prec: 1, Want: "2"},
		19: {V: 1555, Base: PrefixMilli, Prec: 3, Want: "1.56"},
		20: {V: 1550, Base: PrefixMilli, Prec: 2, Want: "1.6"},
		// Rounding close calls.
		21: {V: 999, Base: PrefixMilli, Prec: 3, Want: "999m"},
		22: {V: 999_999, Base: PrefixMilli, Prec: 6, Want: "999.999"},
		23: {V: 9_999_999, Base: PrefixMilli, Prec: 7, Want: "9.999999k"},
		// Rounding events.
		24: {V: 999, Base: PrefixMicro, Prec: 2, Want: "1m"},
		25: {V: 999, Base: PrefixMilli, Prec: 2, Want: "1"},
	}
	s := make([]byte, 24)
	for i, test := range tests {
		if test.Prec == 0 {
			continue // Undeclared
		}
		s = AppendFixed(s[:0], test.V, test.Base, 'f', test.Prec)
		if string(s) != test.Want {
			t.Errorf("case %d: want %s, got %s", i, test.Want, s)
		}
		// Negative equivalent
		s = AppendFixed(s[:0], -test.V, test.Base, 'f', test.Prec)
		if s[0] != '-' || string(s[1:]) != test.Want {
			t.Errorf("case %d: want -%s, got %s", i, test.Want, s)
		}
	}
}

func TestFormatAppend_precorpus(t *testing.T) {
	var tests = []struct {
		V    int64
		Base Prefix
		Prec int
		Want string
	}{}
	s := make([]byte, 24)
	for i, test := range tests {
		if test.Prec == 0 {
			continue // Undeclared
		}
		s = AppendFixed(s[:0], test.V, test.Base, 'f', test.Prec)
		if string(s) != test.Want {
			t.Errorf("case %d: want %s, got %s", i, test.Want, s)
		}
	}
}

func TestNewDimension(t *testing.T) {
	for l := -maxunit; l < maxunit; l += 11 {
		for m := -maxunit; m < maxunit; m += 11 {
			for k := -maxunit; k < maxunit; k += 11 {
				d, err := NewDimension(l, m, k, 4, 5, 6, 7)
				if err != nil {
					t.Fatal("dimension error", err)
				}
				gotl := d.ExpLength()
				if gotl != l {
					t.Errorf("%s:L want %d, got %d", d.String(), l, gotl)
				}
				gotm := d.ExpMass()
				if gotm != m {
					t.Errorf("%s:M want %d, got %d", d.String(), m, gotm)
				}
				gotk := d.ExpTime()
				if gotk != k {
					t.Errorf("%s:K want %d, got %d", d.String(), k, gotk)
				}
				if t.Failed() {
					t.Fatal(t.Name(), "exit early")
				}
			}
		}
	}
}

func TestA(t *testing.T) {
	d, err := NewDimension(1, 2, 3, 4, 5, 6, 6)
	if err != nil {
		panic(err)
	}
	if d.String() != "LM²T³K⁴I⁵J⁶N⁶" {
		t.Fatal("unexpected string positives", d.String())
	}

	d, err = NewDimension(-1, -2, -3, -4, -5, -6, -6)
	if err != nil {
		panic(err)
	}
	if d.String() != "L⁻¹M⁻²T⁻³K⁻⁴I⁻⁵J⁻⁶N⁻⁶" {
		t.Fatal("unexpected string negatives", d.String())
	}
}

func TestIlog10(t *testing.T) {
	var tests = []struct {
		v    int64
		want int
	}{
		{v: 0, want: -1},
		{v: 1, want: 0},
		{v: 5, want: 0},
		{v: 10, want: 1},
		{v: 15, want: 1},
		{v: 50, want: 1},
		{v: 99, want: 1},
		{v: 100, want: 2},
		{v: 999, want: 2},
		{v: 1000, want: 3},
	}
	for i, test := range tests {
		got := ilog10(test.v)
		if got != test.want {
			t.Errorf("case %d: want %d, got %d with %d", i, test.want, got, test.v)
		}
	}
}
