package si

import (
	"math/rand"
	"testing"
)

func TestFormatAppend(t *testing.T) {
	var tests = []struct {
		V     int64
		BaseU Prefix
		Prec  int
		Want  string
	}{
		// Augment from Prefixed to No prefix.
		0: {V: 1000, BaseU: PrefixMilli, Prec: 4, Want: "1"},
		1: {V: 10000, BaseU: PrefixMilli, Prec: 4, Want: "10"},
		2: {V: 100000, BaseU: PrefixMilli, Prec: 4, Want: "100"},
		// Print within prefix range
		3: {V: 1, BaseU: PrefixMilli, Prec: 4, Want: "1m"},
		4: {V: 10, BaseU: PrefixMilli, Prec: 4, Want: "10m"},
		5: {V: 100, BaseU: PrefixMilli, Prec: 4, Want: "100m"},
		// Augment from low Prefixed to higher prefix.
		6: {V: 1_000_000, BaseU: PrefixMilli, Prec: 4, Want: "1k"},
		7: {V: 10_000_000, BaseU: PrefixMilli, Prec: 4, Want: "10k"},
		8: {V: 100_000_000, BaseU: PrefixMilli, Prec: 4, Want: "100k"},
		// Augment Prefix to No Prefix with decimals.
		9:  {V: 1234, BaseU: PrefixMilli, Prec: 4, Want: "1.234"},
		10: {V: 12344, BaseU: PrefixMilli, Prec: 4, Want: "12.34"},
		11: {V: 123444, BaseU: PrefixMilli, Prec: 4, Want: "123.4"},
		// Augment low to high prefix with decimals.
		12: {V: 1_234_444, BaseU: PrefixMilli, Prec: 4, Want: "1.234k"},
		13: {V: 12_344_444, BaseU: PrefixMilli, Prec: 4, Want: "12.34k"},
		14: {V: 123_444_444, BaseU: PrefixMilli, Prec: 4, Want: "123.4k"},
		// Augment low to high prefix with decimal chop-off.
		15: {V: 1_234_567, BaseU: PrefixMilli, Prec: 1, Want: "1k"},
		16: {V: 12_345_678, BaseU: PrefixMilli, Prec: 2, Want: "12k"},
		17: {V: 123_456_789, BaseU: PrefixMilli, Prec: 3, Want: "123k"},
		// Rounding simple.
		18: {V: 1500, BaseU: PrefixMilli, Prec: 1, Want: "2"},
		19: {V: 1555, BaseU: PrefixMilli, Prec: 3, Want: "1.56"},
		20: {V: 1550, BaseU: PrefixMilli, Prec: 2, Want: "1.6"},
		// Rounding close calls.
		21: {V: 999, BaseU: PrefixMilli, Prec: 3, Want: "999m"},
		22: {V: 999_999, BaseU: PrefixMilli, Prec: 6, Want: "999.999"},
		23: {V: 9_999_999, BaseU: PrefixMilli, Prec: 7, Want: "9.999999k"},
		// Normal rounding events.
		24: {V: 12345, BaseU: PrefixMilli, Prec: 4, Want: "12.35"},
		25: {V: 1500, BaseU: PrefixMilli, Prec: 1, Want: "2"},
		// 26: {V: 15, Base: PrefixMilli, Prec: 2, Want: "15m"},
		// Extraordinary base-crossing rounding events.
		27: {V: 999_999, BaseU: PrefixMicro, Prec: 2, Want: "1"},
		28: {V: 999_999, BaseU: PrefixMilli, Prec: 2, Want: "1k"},
	}
	s := make([]byte, 24)
	for i, test := range tests {
		if test.Prec == 0 {
			continue // Undeclared or commented.
		}
		s = AppendFixed(s[:0], test.V, test.BaseU, 'f', test.Prec)
		if string(s) != test.Want {
			t.Errorf("case %d: want %s, got %s", i, test.Want, s)
		}
		// Negative equivalent
		s = AppendFixed(s[:0], -test.V, test.BaseU, 'f', test.Prec)
		if s[0] != '-' || string(s[1:]) != test.Want {
			t.Errorf("case %d: want -%s, got %s", i, test.Want, s)
		}
	}
}

func TestFormatAppend_precorpus(t *testing.T) {
	var tests = []struct {
		V     int64
		BaseU Prefix
		Prec  int
		Want  string
	}{
		// {V: 15, Base: PrefixMilli, Prec: 1, Want: "123.5"},
	}
	s := make([]byte, 24)
	for i, test := range tests {
		if test.Prec == 0 {
			continue // Undeclared
		}
		s = AppendFixed(s[:0], test.V, test.BaseU, 'f', test.Prec)
		if string(s) != test.Want {
			t.Errorf("case %d: want %s, got %s", i, test.Want, s)
		}
	}
}

func TestNewDimension(t *testing.T) {
	const testMaxUnit = 127
	for l := -testMaxUnit; l < testMaxUnit; l += 11 {
		for m := -testMaxUnit; m < testMaxUnit; m += 11 {
			for k := -testMaxUnit; k < testMaxUnit; k += 11 {
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

func TestFixedToFloat(t *testing.T) {
	var tests = []struct {
		V     int64
		BaseU Prefix
		Want  float64
	}{
		// Basic conversions.
		0: {V: 1, BaseU: PrefixNone, Want: 1.0},
		1: {V: 0, BaseU: PrefixNone, Want: 0.0},
		2: {V: 1000, BaseU: PrefixNone, Want: 1000.0},
		3: {V: -1000, BaseU: PrefixNone, Want: -1000.0},
		4: {V: 1_000_000, BaseU: PrefixNone, Want: 1_000_000.0},
		// Milli base (10^-3).
		5:  {V: 1, BaseU: PrefixMilli, Want: 0.001},
		6:  {V: 1000, BaseU: PrefixMilli, Want: 1.0},
		7:  {V: 1_000_000, BaseU: PrefixMilli, Want: 1000.0},
		8:  {V: 2500, BaseU: PrefixMilli, Want: 2.5},
		9:  {V: 123_456, BaseU: PrefixMilli, Want: 123.456},
		10: {V: -1000, BaseU: PrefixMilli, Want: -1.0},
		// Kilo base (10^3).
		11: {V: 1, BaseU: PrefixKilo, Want: 1000.0},
		12: {V: 500, BaseU: PrefixKilo, Want: 500_000.0},
		13: {V: 1000, BaseU: PrefixKilo, Want: 1_000_000.0},
		14: {V: -2, BaseU: PrefixKilo, Want: -2000.0},
		// Micro base (10^-6).
		15: {V: 1, BaseU: PrefixMicro, Want: 0.000001},
		16: {V: 1_000_000, BaseU: PrefixMicro, Want: 1.0},
		17: {V: 500_000, BaseU: PrefixMicro, Want: 0.5},
		// Mega base (10^6).
		18: {V: 1, BaseU: PrefixMega, Want: 1_000_000.0},
		19: {V: 2, BaseU: PrefixMega, Want: 2_000_000.0},
		20: {V: 1000, BaseU: PrefixMega, Want: 1_000_000_000.0},
		// Giga base (10^9).
		21: {V: 1, BaseU: PrefixGiga, Want: 1_000_000_000.0},
		22: {V: 5, BaseU: PrefixGiga, Want: 5_000_000_000.0},
		// Extreme cases.
		23: {V: 1, BaseU: PrefixAtto, Want: 1e-18},
		24: {V: 1, BaseU: PrefixPeta, Want: 1e15},
		25: {V: 1, BaseU: PrefixExa, Want: 1e18},
	}
	for i, test := range tests {
		got := FixedToFloat(test.V, test.BaseU)
		if got != test.Want {
			t.Errorf("case %d: FixedToFloat(%d, %d) = %g, want %g", i, test.V, test.BaseU, got, test.Want)
		}
	}
}

func TestParseFixed(t *testing.T) {
	var tests = []struct {
		S     string
		BaseU Prefix
		Want  int64
	}{
		// Usual cases
		0: {S: "1k", BaseU: PrefixMilli, Want: 1_000_000},
		1: {S: "2M", BaseU: PrefixMilli, Want: 2_000_000_000},
		2: {S: "1.0004M", BaseU: PrefixMilli, Want: 1_000_400_000},
		// Special preceding dot cases.
		3: {S: "0.10k", BaseU: PrefixMilli, Want: 100_000},
		4: {S: "0.010k", BaseU: PrefixMilli, Want: 10_000},
		5: {S: "0.0010k", BaseU: PrefixMilli, Want: 1_000},
		// Exponent notation with lowercase 'e'.
		6:  {S: "2e5m", BaseU: PrefixMilli, Want: 200_000},
		7:  {S: "3e2k", BaseU: PrefixMilli, Want: 300_000_000},
		8:  {S: "1e3", BaseU: PrefixMilli, Want: 1_000_000},
		9:  {S: "5e1m", BaseU: PrefixMilli, Want: 50},
		10: {S: "1e0k", BaseU: PrefixMilli, Want: 1_000_000},
		// Exponent notation with uppercase 'E'.
		11: {S: "2E5m", BaseU: PrefixMilli, Want: 200_000},
		12: {S: "3E2k", BaseU: PrefixMilli, Want: 300_000_000},
		13: {S: "1E3", BaseU: PrefixMilli, Want: 1_000_000},
		// Negative exponents.
		14: {S: "5e-2k", BaseU: PrefixMilli, Want: 50_000},
		15: {S: "2e-3m", BaseU: PrefixMilli, Want: 0}, // 0.002 milli -> rounds to 0
		16: {S: "123e-1k", BaseU: PrefixMilli, Want: 12_300_000},
		17: {S: "5E-1m", BaseU: PrefixMilli, Want: 1}, // 0.5 milli -> rounds to 1
		// Positive exponents with explicit '+' sign.
		18: {S: "2e+5m", BaseU: PrefixMilli, Want: 200_000},
		19: {S: "3E+2k", BaseU: PrefixMilli, Want: 300_000_000},
		20: {S: "1e+0k", BaseU: PrefixMilli, Want: 1_000_000},
		// Exponent notation with decimal points.
		21: {S: "1.5e3k", BaseU: PrefixMilli, Want: 1_500_000_000},
		22: {S: "2.5e2m", BaseU: PrefixMilli, Want: 250},
		23: {S: "1.23e4", BaseU: PrefixMilli, Want: 12_300_000},
		24: {S: "4.56e-1k", BaseU: PrefixMilli, Want: 456_000},
		// Different base units.
		25: {S: "1e6", BaseU: PrefixNone, Want: 1_000_000},
		26: {S: "2e3k", BaseU: PrefixNone, Want: 2_000_000},
		27: {S: "5e-3M", BaseU: PrefixNone, Want: 5_000},
		28: {S: "1e11m", BaseU: PrefixKilo, Want: 100_000}, // 1e11 milli = 1e8 / 1e3 kilo = 1e5 = 100k
		// Multi-digit exponents.
		29: {S: "1e10", BaseU: PrefixMilli, Want: 10_000_000_000_000},
		30: {S: "2e12m", BaseU: PrefixMilli, Want: 2_000_000_000_000},
		// Peta and Exa prefixes (uppercase 'P' and 'E' without exponent notation).
		31: {S: "2P", BaseU: PrefixMilli, Want: 2_000_000_000_000_000_000},
		32: {S: "1P", BaseU: PrefixNone, Want: 1_000_000_000_000_000},
		33: {S: "3P", BaseU: PrefixKilo, Want: 3_000_000_000_000},
		34: {S: "1E", BaseU: PrefixKilo, Want: 1_000_000_000_000_000},
		35: {S: "9E", BaseU: PrefixMega, Want: 9_000_000_000_000},
	}
	for i, test := range tests {
		if test.S == "" {
			continue // Commented line.
		}
		v, n, err := ParseFixed(test.S, test.BaseU)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(test.S) {
			t.Errorf("case %d: bytes read mismatch, got %d want %d", i, n, len(test.S))
		}
		if v != test.Want {
			t.Errorf("case %d: got %d, want %d from %q with baseUnits=%d", i, v, test.Want, test.S, test.BaseU)
		}
	}
}

func TestFormatParseLoop(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	var buf [24]byte
	for i := 0; i < 10000; i++ {
		const baseUnits = PrefixAtto
		v := rng.Int63() / 16
		digits := ilog10(v) + 1
		if rng.Int()%2 != 0 {
			v = -v // Do negative case.
		}
		s := AppendFixed(buf[:0], v, baseUnits, 'f', digits)
		if v >= 0 && v%16 == 0 {
			s = append([]byte{'+'}, s...) // Do leading positive number case in 1/16 of positive cases.
		}
		got, n, err := ParseFixed(string(s), baseUnits)
		if err != nil {
			t.Fatalf("%s: %q", err, s)
		} else if n != len(s) {
			t.Fatalf("mismatch number of bytes read got %d, want %d for %q", n, len(s), s)
		}
		if got != v {
			t.Errorf("format-parse loop failed for got %d, want %d for %q", got, v, s)
		}
	}
}

func TestStringStuff(t *testing.T) {
	si, _ := NewDimensionFormatter(DefaultDimensionFormatterConfig())
	d, _ := NewDimension(1, 2, 3, 4, 5, 6, 7)
	if si.StringDim(d) != "m·kg²·s³·K⁴·A⁵·cd⁶·mol⁷" {
		t.Error("bad format", si.StringDim(d))
	}
}

func TestAppendFixedErrors(t *testing.T) {
	var tests = []struct {
		V     int64
		BaseU Prefix
		Prec  int
	}{
		// Invalid prec.
		0: {V: 1234, Prec: -1},
		1: {V: 1234, Prec: 0},
		2: {V: 1234, Prec: 22},
		// Invalid Prefix.
		3: {V: 1234, Prec: 1, BaseU: 1},
		4: {V: 1234, Prec: 1, BaseU: PrefixExa + 3},
		5: {V: 1234, Prec: 1, BaseU: PrefixAtto - 3},
		// Exceed base units upwards (beyond PrefixExa).
		6: {V: 1234, Prec: 3, BaseU: PrefixExa},
		7: {V: 1234567, Prec: 1, BaseU: PrefixPeta},
		8: {V: 1234567, Prec: 3, BaseU: PrefixPeta},
	}
	var buf [64]byte
	for i, test := range tests {
		s := AppendFixed(buf[:0], test.V, test.BaseU, 'f', test.Prec)
		if s[0] != '<' {
			t.Errorf("case %d: expected error, got %q", i, s)
		}
		s[0] = '0' // reset to avoid flukes.
	}
}

func TestParseFixedErrors(t *testing.T) {
	var tests = []struct {
		S     string
		BaseU Prefix
	}{
		// Bad dots.
		0: {S: "1..234"},
		1: {S: ".1234."},
		2: {S: "1.2.34"},
		// Bad Minus.
		3: {S: "--1234"},
		4: {S: "1-234"},
		5: {S: "-1234-"},
		// Bad Plus.
		6: {S: "1+234"},
		7: {S: "++1234"},
		8: {S: "+1+234"},
		// Bad Plus/minus
		9:  {S: "-+1234"},
		10: {S: "+-1234"},
		11: {S: "+1234-"},
		// Digit overflow.
		12: {S: "12345678901234567890"},  // 20 digits.
		13: {S: "1.2345678901234567890"}, // 20 digits with dot.
		14: {S: "12345678901234567890."},
		// Base unit overflow.
		15: {S: "12345678901234567", BaseU: PrefixMilli},
		16: {S: "12345678901234567.000", BaseU: PrefixMilli},
		17: {S: "12345678901234", BaseU: PrefixMicro},
		// Bad exponent notation.
		18: {S: "2e"},    // Lowercase 'e' at end is unknown prefix (not Exa).
		19: {S: "2e-"},   // Missing exponent digits after sign.
		20: {S: "2e+"},   // Missing exponent digits after sign.
		21: {S: "2E-"},   // Missing exponent digits after sign (uppercase).
		22: {S: "2E+"},   // Missing exponent digits after sign (uppercase).
		23: {S: "2e--3"}, // Double negative in exponent.
		24: {S: "2e++3"}, // Double plus in exponent.
		25: {S: "2e+-3"}, // Plus and minus in exponent.
		26: {S: "2e-+3"}, // Minus and plus in exponent.
	}
	for i, test := range tests {
		v, n, err := ParseFixed(test.S, test.BaseU)
		if err == nil {
			t.Fatalf("case %d: expected error, got %d from %q", i, v, test.S)
		}
		if n != 0 {
			t.Errorf("case %d: expected no bytes read on error, got %d from %q", i, n, test.S)
		} else if v != 0 {
			t.Errorf("case %d: expected zero output value, got %d from %q", i, v, test.S)
		}
	}
}

func TestDimensionAlgebra(t *testing.T) {
	// unit, _ := NewDimension(1, 1, 1, 1, 1, 1, 1)
	var units [7]Dimension
	for i := range units {
		s := make([]int, 7)
		s[i] = 1
		units[i] = newdim(s)
	}
	var d Dimension
	// Test multiplication.
	for i, udim := range units {
		k, err := MulDim(d, udim)
		if err != nil {
			t.Fatal(err)
		}
		for j, exp := range k.Exponents() {
			if j == i && exp != 1 {
				t.Error("bad dimension add")
			} else if j != i && exp != 0 {
				t.Error("out of band dimension added")
			}
		}
	}
	// Test Division.
	for i, udim := range units {
		k, err := DivDim(d, udim)
		if err != nil {
			t.Fatal(err)
		}
		for j, exp := range k.Exponents() {
			if j == i && exp != -1 {
				t.Error("bad dimension add")
			} else if j != i && exp != 0 {
				t.Error("out of band dimension added")
			}
		}
		k = k.Inv()
		for j, exp := range k.Exponents() {
			if j == i && exp != 1 {
				t.Error("bad dimension add")
			} else if j != i && exp != 0 {
				t.Error("out of band dimension added")
			}
		}
	}
}

func newdim(dims []int) Dimension {
	d, err := NewDimension(dims[0], dims[1], dims[2], dims[3], dims[4], dims[5], dims[6])
	if err != nil {
		panic(err)
	}
	return d
}
