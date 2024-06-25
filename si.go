package si

import (
	"errors"
	"math"
	"strconv"
	"unicode/utf8"
)

type dimint = int8

const (
	maxunit = math.MaxInt8
)

var errDimOOB = errors.New("dimension exceeds storage space (-127..127)")

// NewDimension creates a new dimension from the given exponents.
func NewDimension(Length, Mass, Time, Temperature, ElectricCurrent, Luminosity, Amount int) (Dimension, error) {
	if isDimOOB(Length) || isDimOOB(Mass) || isDimOOB(Time) ||
		isDimOOB(Temperature) || isDimOOB(ElectricCurrent) ||
		isDimOOB(Luminosity) || isDimOOB(Amount) {
		return Dimension{}, errDimOOB
	}
	return Dimension{
		dims: [7]dimint{
			0: dimint(Length),
			1: dimint(Mass),
			2: dimint(Time),
			3: dimint(Temperature),
			4: dimint(ElectricCurrent),
			5: dimint(Luminosity),
			6: dimint(Amount),
		},
	}, nil
}

func isDimOOB(dim int) bool {
	return dim > maxunit || dim < -maxunit
}

// Dimension represents the dimensions of a physical quantity.
type Dimension struct {
	// dims contains 7 int8's representing the exponent of primitive dimensions.
	// The ordering follows the result of Exponents method result.
	dims [7]dimint
}

const negexp = '⁻'

var exprune = [10]rune{
	0: '⁰',
	1: '¹',
	2: '²',
	3: '³',
	4: '⁴',
	5: '⁵',
	6: '⁶',
	7: '⁷',
	8: '⁸',
	9: '⁹',
}

var (
	abstractDimFormatter, _ = NewDimensionFormatter(AbstractDimensionFormatterConfig())
)

// AbstractDimensionFormatConfig returns the abstract unit formatting configuration.
func AbstractDimensionFormatterConfig() DimensionFormatterConfig {
	return DimensionFormatterConfig{
		Length:      "L",
		Mass:        "M",
		Time:        "T",
		Temperature: "K",
		Current:     "I",
		Luminosity:  "J",
		Amount:      "N",
	}
}

// DefaultDimensionFormatter returns the SI formatter config.
func DefaultDimensionFormatterConfig() DimensionFormatterConfig {
	return DimensionFormatterConfig{
		Length:      "m",
		Mass:        "kg",
		Time:        "s",
		Temperature: "K",
		Current:     "A",
		Luminosity:  "cd",
		Amount:      "mol",
		Separator:   "·",
	}
}

// DimensionFormatter is an arrangement of unit representations.
type DimensionFormatter struct {
	fmts [7]string
	sep  string
}

// DimensionFormatterConfig specifies how the DimensionFormatter will
// behave during formatting calls on the Dimension type.
type DimensionFormatterConfig struct {
	Length      string
	Mass        string
	Time        string
	Temperature string
	Current     string
	Luminosity  string
	Amount      string

	// Unit separator.
	Separator string
}

// NewDimensionFormatter creates a new dimension formatter.
func NewDimensionFormatter(cfg DimensionFormatterConfig) (*DimensionFormatter, error) {
	if cfg.Length == "" || cfg.Mass == "" || cfg.Time == "" || cfg.Temperature == "" || cfg.Current == "" || cfg.Luminosity == "" || cfg.Amount == "" {
		return nil, errors.New("empty format string")
	}

	return &DimensionFormatter{
		sep: cfg.Separator,
		fmts: [7]string{
			0: cfg.Length,
			1: cfg.Mass,
			2: cfg.Time,
			3: cfg.Temperature,
			4: cfg.Current,
			5: cfg.Luminosity,
			6: cfg.Amount,
		},
	}, nil
}

// String returns a human readable representation of the DimensionFormatter.
func (df *DimensionFormatter) String() string {
	b := make([]byte, 0, 8*4) // 32 stores SI perfectly.
	for i := range df.fmts {
		b = append(b, abstractDimFormatter.fmts[i]...)
		b = append(b, ':')
		b = append(b, df.fmts[i]...)
		if i != len(df.fmts)-1 {
			b = append(b, ' ')
		}
	}
	return string(b)
}

// StringDim returns the string representation of the dimension with df's formatting directive.
func (df *DimensionFormatter) StringDim(dim Dimension) string {
	if dim.IsDimensionless() {
		return ""
	}
	return string(df.AppendFormat(make([]byte, 0, df.sizeofFormat(dim)), dim))
}

// sizeofFormat returns exact size of the printed string of dim with df formatting.
func (df *DimensionFormatter) sizeofFormat(dim Dimension) int {
	sizeof := 0
	var printed bool
	for i, exp := range dim.dims {
		if exp == 0 {
			// Exponent not printed.
			continue
		}
		if printed {
			sizeof += len(df.sep)
		}
		printed = true
		// Size of unit string added.
		sizeof += len(df.fmts[i])
		// Size of exponent in bytes, exponents are 2 bytes in length.
		if exp < 0 {
			sizeof += 2
		}
		digits := ilog10(int64(exp)) + 1
		sizeof += digits * 2
	}
	return sizeof
}

// AppendFormat formats a dimension with
func (df *DimensionFormatter) AppendFormat(b []byte, dim Dimension) []byte {
	// Size of this buffer should fit -32767 (MaxInt16)
	if dim.IsDimensionless() {
		return b
	}
	var buf [8]byte
	var lastPrinted bool
	for i := range df.fmts {
		dim := dim.dims[i]
		if dim == 0 {
			continue
		}
		if lastPrinted {
			b = append(b, df.sep...)
		}
		lastPrinted = true
		b = append(b, df.fmts[i]...)
		if dim == 1 {
			continue
		}
		numbuf := strconv.AppendInt(buf[:0], int64(dim), 10)
		if numbuf[0] == '-' {
			b = utf8.AppendRune(b, negexp)
			numbuf = numbuf[1:]
		}
		for i := 0; i < len(numbuf); i++ {
			offset := numbuf[i] - '0'
			if offset > 9 {
				panic("invalid char")
			}
			b = utf8.AppendRune(b, exprune[offset])
		}
	}
	return b
}

// String returns a human readable representation of the dimension using abstract unit letters (LMTKIJN).
func (d Dimension) String() string {
	return abstractDimFormatter.StringDim(d)
}

// IsDimensionless returns true if d is dimensionless, that is to say all dimension exponents are zero.
func (d Dimension) IsDimensionless() bool { return d.dims == ([7]dimint{}) }

// ExpLength returns the exponent of the length dimension of d.
func (d Dimension) ExpLength() int { return int(d.dims[0]) }

// ExpMass returns the exponent of the mass dimension of d.
func (d Dimension) ExpMass() int { return int(d.dims[1]) }

// ExpTime returns the exponent of the time dimension of d.
func (d Dimension) ExpTime() int { return int(d.dims[2]) }

// ExpTemperature returns the exponent of the temperature dimension of d.
func (d Dimension) ExpTemperature() int { return int(d.dims[3]) }

// ExpCurrent returns the exponent of the current dimension of d.
func (d Dimension) ExpCurrent() int { return int(d.dims[4]) }

// ExpLuminosity returns the exponent of the luminosity dimension of d.
func (d Dimension) ExpLuminosity() int { return int(d.dims[5]) }

// ExpAmount returns the exponent of the amount dimension of d.
func (d Dimension) ExpAmount() int { return int(d.dims[6]) }

// Exponents returns the exponents of the 7 dimensions as an array. The ordering is:
//  0. Distance dimension (L)
//  1. Mass dimension (M)
//  2. Time dimension (T)
//  3. Temperature dimension (K)
//  4. Electric current dimension (I)
//  5. Luminosity intensity dimension (J)
//  6. Amount or quantity dimension (N)
func (d Dimension) Exponents() (LMTKIJN [7]int) {
	for i := range LMTKIJN {
		LMTKIJN[i] = int(d.dims[i])
	}
	return LMTKIJN
}

// Inv inverts the dimension by multiplying all dimension exponents by -1.
func (d Dimension) Inv() Dimension {
	inv := d
	for i := range inv.dims {
		inv.dims[i] *= -1
	}
	return inv
}

// MulDim returns the dimension obtained from a*b.
// It returns an error if result dimension exceeds storage.
func MulDim(a, b Dimension) (Dimension, error) {
	L := a.ExpLength() + b.ExpLength()
	M := a.ExpMass() + b.ExpMass()
	T := a.ExpTime() + b.ExpTime()
	K := a.ExpTemperature() + b.ExpTemperature()
	I := a.ExpCurrent() + b.ExpCurrent()
	J := a.ExpLuminosity() + b.ExpLuminosity()
	N := a.ExpAmount() + b.ExpAmount()
	return NewDimension(L, M, T, K, I, J, N)
}

// DivDim returns the dimension obtained from a/b.
// It returns an error if result dimension exceeds storage.
func DivDim(a, b Dimension) (Dimension, error) {
	return MulDim(a, b.Inv())
}

// Prefix represents a unit prefix used to specify the magnitude of a quantity.
// i.e: PrefixKilo corresponds to 'k' character used to denote a multiplier of 1000 to the unit it is prefixed to.
type Prefix int8

// Package unit prefix definitions.
const (
	PrefixAtto Prefix = -18 + iota*3
	PrefixFemto
	PrefixPico
	PrefixNano
	PrefixMicro
	PrefixMilli
	PrefixNone
	PrefixKilo
	PrefixMega
	PrefixGiga
	PrefixTera
	PrefixExa
)

// RuneToPrefix interprets the argument rune as an SI prefix.
// It does not parse PrefixNone.
func RuneToPrefix(r rune) (pfx Prefix, err error) {
	switch r {
	case 'a':
		pfx = PrefixAtto
	case 'f':
		pfx = PrefixFemto
	case 'p':
		pfx = PrefixPico
	case 'n':
		pfx = PrefixNano
	case 'u', 'μ':
		pfx = PrefixMicro
	case 'm':
		pfx = PrefixMilli
	case 'k':
		pfx = PrefixKilo
	case 'M':
		pfx = PrefixMega
	case 'G':
		pfx = PrefixGiga
	case 'T':
		pfx = PrefixTera
	case 'E':
		pfx = PrefixExa
	default:
		err = errUnknownPrefix
	}
	return pfx, err
}

// IsValid checks if the prefix is one of the supported standard SI prefixes or the zero base prefix.
func (p Prefix) IsValid() bool {
	return p == PrefixNone || p.Character() != ' '
}

// String returns a human readable representation of the Prefix of string type.
// Returns a error message string if Prefix is undefined. Guarateed to return non-zero string.
func (p Prefix) String() string {
	if p == PrefixMicro {
		return "μ"
	}
	const pfxTable = "a!!f!!p!!n!!u!!m!! !!k!!M!!G!!T!!E"
	offset := int(p - PrefixAtto)
	if offset < 0 || offset >= len(pfxTable) || pfxTable[offset] == '!' {
		return "<si!invalid Prefix>"
	}
	return pfxTable[offset : offset+1]
}

// Character returns the single character SI representation of the unit prefix.
// If not representable or invalid returns space caracter ' '.
func (p Prefix) Character() (s rune) {
	if p == PrefixMicro {
		return 'μ'
	}
	s = rune(p.String()[0])
	if s == '<' {
		s = ' '
	}
	return s
}

// AppendFixed formats a fixed-point number with a given magnitude base units and
// appends it's representation to the argument buffer.
//
//	"123.456k" for value=123456, baseUnits=PrefixNone, prec=6
//	"123k" for value=123456, baseUnits=PrefixNone, prec=3
func AppendFixed(b []byte, value int64, baseUnits Prefix, fmt byte, prec int) []byte {
	switch {
	case fmt != 'f':
		return append(b, "<si!INVALID FMT>"...)
	case prec <= 0:
		return append(b, "<si!LESS-EQ-ZERO PREC>"...)
	case !baseUnits.IsValid():
		return append(b, "<si!BAD BASE>"...)
	case value == 0:
		return append(b, '0')
	}

	isNegative := value < 0
	if isNegative {
		value = -value
	}
	v64 := int64(value)
	log10 := ilog10(v64)

	log10mod3 := log10 % 3
	frontDigits := (log10mod3) + 1
	backDigits := log10 - log10mod3
	if log10 < backDigits {
		frontDigits = 0
	}

	// We now prepare the value by trimming digits after the precision cutoff.
	// We need to trim when excess > 0.
	// TODO: Decide on whether to keep forced 3-sigfig formatting when only at 3 digits corresponding to `backDigits != 0` condition.
	excess := backDigits + frontDigits - prec
	if excess > 0 && backDigits != 0 {
		x := v64 % powerOf10[excess]
		rlim := iLogRoundTable[excess]
		roundUp := x >= rlim
		v64 /= powerOf10[excess]
		v64 += int64(b2i(roundUp))
		newIlog10 := ilog10(v64) + excess
		if newIlog10 != log10 {
			// Rounding where frontdigits overflow.
			log10 = newIlog10
			log10mod3 = log10 % 3
			frontDigits = log10mod3 + 1
			backDigits = log10 - log10mod3
			if log10 < backDigits {
				frontDigits = 0
			}
		}
	}

	if isNegative {
		b = append(b, '-')
	}

	var buf [20]byte
	prevlen := len(b)
	b = strconv.AppendInt(b, v64, 10)
	last := append(buf[:0], b[prevlen+frontDigits:]...)
	b = b[:prevlen+frontDigits]

	for i := range last {
		// Only print if has non-zero part.
		if last[i] != '0' {
			b = append(b[:prevlen+frontDigits], '.')
			b = append(b, last...)
			break
		}
	}

	// Calculate new base.
	baseUnits += Prefix(log10 - log10mod3)
	if baseUnits != PrefixNone {
		b = append(b, baseUnits.String()...)
	}
	return b
}

// ParseFixed parses a decimal point representation with or without unit prefix
// and converts it to a fixed point representation with `baseUnits` as the base units.
func ParseFixed(s string, baseUnits Prefix) (value int64, readBytes int, err error) {
	var buf [20]byte
	// s indices.
	var dotPos, wholeEnd, bufPtr int = -1, 0, 0
	var seenPlus bool
	var d decimal
CHARLOOP:
	for wholeEnd < len(s) {
		c := s[wholeEnd]
		if '0' <= c && c <= '9' {
			if bufPtr >= len(buf) {
				err = errOverflowsInt64
				break CHARLOOP
			} else if bufPtr == 0 && dotPos < 0 && c == '0' {
				// Skip zeros preceding decimal point.
				wholeEnd++
				continue
			}
			buf[bufPtr] = c
			bufPtr++
			wholeEnd++
			if dotPos >= 0 {
				// Digits correspond to decimal part, so subtract from exp.
				d.exp--
			}
			continue
		}
		switch c {
		case '.':
			if dotPos >= 0 {
				err = errDotDot
				break CHARLOOP
			}
			dotPos = wholeEnd
		case '+':
			if seenPlus {
				err = errPlusPlus
				break CHARLOOP
			} else if d.neg {
				err = errPlusMinus
				break CHARLOOP
			} else if wholeEnd != 0 {
				err = errNaN
				break CHARLOOP
			}
			seenPlus = true
		case '-':
			if d.neg {
				err = errMinusMinus
				break CHARLOOP
			} else if seenPlus {
				err = errPlusPlus
				break CHARLOOP
			} else if wholeEnd != 0 {
				err = errNaN
				break CHARLOOP
			}
			d.neg = true
		default:
			break CHARLOOP
		}
		wholeEnd++
	}
	if err != nil {
		return 0, 0, err
	}
	readBytes = wholeEnd
	var incomingPrefix Prefix
	if wholeEnd < len(s) {
		r, n := utf8.DecodeRuneInString(s[readBytes:])
		incomingPrefix, err = RuneToPrefix(r)
		if err != nil {
			return 0, 0, err
		}
		readBytes += n
	}
	// Calculate exponent modifier from decimal point.
	// Where bufPtr is length of number, dotPos is position of decimal w.r.t start.
	//  xxx.xxxxxx gives dotPos=3, bufPtr=3+6 -> exp=-6
	d.base, err = strconv.ParseUint(string(buf[:bufPtr]), 10, 64)
	if err != nil {
		return 0, 0, err
	}

	v, overflow := dtoi(d, int(incomingPrefix-baseUnits))
	if overflow {
		return 0, 0, errOverflowsInt64
	}
	return v, readBytes, nil
}

// ilog10 returns the integer logarithm base 10 of v, which
// can be interpreted as the quanity of digits in the number in base 10 minus one.
func ilog10(v int64) int {
	for i, l := range powerOf10 {
		if v < l {
			return i - 1
		}
	}
	return len(powerOf10)
}

var powerOf10 = [...]int64{
	1,
	10,
	100,
	1_000,
	10_000,
	100_000,
	1_000_000,
	10_000_000,
	100_000_000,
	1_000_000_000,
	10_000_000_000,
	100_000_000_000,
	1_000_000_000_000,
	10_000_000_000_000,
	100_000_000_000_000,
	1_000_000_000_000_000,
	10_000_000_000_000_000,
	100_000_000_000_000_000,
	1_000_000_000_000_000_000,
}

var iLogRoundTable = [...]int64{
	0,
	5,
	50,
	500,
	5_000,
	50_000,
	500_000,
	5_000_000,
	50_000_000,
	500_000_000,
	5_000_000_000,
	50_000_000_000,
	500_000_000_000,
	5_000_000_000_000,
	50_000_000_000_000,
	500_000_000_000_000,
	5_000_000_000_000_000,
	50_000_000_000_000_000,
	500_000_000_000_000_000,
	5_000_000_000_000_000_000,
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

/*
Logic below ripped directly from https://github.com/periph/conn/blob/main/physic/units.go
Modified to avoid ongoing heap allocations and export ParseError for user convenience.
*/

// Decimal is the representation of decimal number.
type decimal struct {
	// base hold the significant digits.
	base uint64
	// exponent is the left or right decimal shift. (powers of ten).
	exp int
	// neg it true if the number is negative.
	neg bool
}
type ParseError struct {
	error
}

func pe(s string) *ParseError {
	return &ParseError{error: errors.New(s)}
}

func (pe *ParseError) Unwrap() error { return pe.error }

var (
	errPlusMinus              = pe("contains both plus and minus")
	errMinusMinus             = pe("contains multiple minus symbols")
	errNaN                    = pe("not a number")
	errPlusPlus               = pe("contains multiple plus symbols")
	errDotDot                 = pe("contains multiple decimal points")
	errOverflowsInt64         = pe("exceeds maximum")
	errOverflowsInt64Negative = pe("exceeds minimum")
	errUnknownPrefix          = pe("unknown SI prefix")
)

// Converts from decimal to int64.
//
// Scale is combined with the decimal exponent to maximise the resolution and is
// in powers of ten.
//
// Returns true if the value overflowed.
func dtoi(d decimal, scale int) (int64, bool) {
	// Get the total magnitude of the number.
	// a^x * b^y = a*b^(x+y) since scale is of the order unity this becomes
	// 1^x * b^y = b^(x+y).
	// mag must be positive to use as index in to powerOf10 array.
	u := d.base
	mag := d.exp + scale
	if mag < 0 {
		mag = -mag
	}
	var n int64
	if mag > 18 {
		return 0, true
	}
	// Divide is = 10^(-mag)
	switch {
	case d.exp+scale < 0:
		u = (u + uint64(powerOf10[mag])/2) / uint64(powerOf10[mag])
	case mag == 0:
		if u > math.MaxInt64 {
			return 0, true
		}
	default:
		check := u * uint64(powerOf10[mag])
		if check/uint64(powerOf10[mag]) != u || check > math.MaxInt64 {
			return 0, true
		}
		u *= uint64(powerOf10[mag])
	}

	n = int64(u)
	if d.neg {
		n = -n
	}
	return n, false
}
