package si

import (
	"errors"
	"strconv"
	"unicode/utf8"
)

const (
	maxunit = 200
)

// NewDimension creates a new dimension from the given exponents.
func NewDimension(length, mass, time, temperature, electricCurrent, luminosity, amount int) (Dimension, error) {
	if isDimOOB(length) || isDimOOB(mass) || isDimOOB(time) ||
		isDimOOB(temperature) || isDimOOB(electricCurrent) ||
		isDimOOB(luminosity) || isDimOOB(amount) {
		return Dimension{}, errors.New("overflow of dimension storage")
	}
	return Dimension{
		dims: [7]int16{
			int16(length), int16(mass),
			int16(time), int16(temperature),
			int16(electricCurrent), int16(luminosity),
			int16(amount),
		},
	}, nil
}

func isDimOOB(dim int) bool {
	return dim > maxunit || dim < -maxunit
}

// Dimension represents the dimensions of a physical quantity.
type Dimension struct {
	// dims contains 7 int16's representing the exponent of primitive dimensions.
	// The ordering follows the result of Exponents method result.
	dims [7]int16
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

// String returns a human readable representation of the dimension.
func (d Dimension) String() string {
	if d.IsDimensionless() {
		return ""
	}
	s := make([]byte, 0, 8)
	return string(d.appendf(s))
}

func (d Dimension) appendf(b []byte) []byte {
	app := func(b []byte, char byte, dim int) []byte {
		if dim == 0 {
			return b
		}
		b = append(b, char)
		if dim == 1 {
			return b
		}

		var buf [20]byte
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
		return b
	}
	b = app(b, 'L', d.ExpLength())
	b = app(b, 'M', d.ExpMass())
	b = app(b, 'T', d.ExpTime())
	b = app(b, 'K', d.ExpTemperature())
	b = app(b, 'I', d.ExpCurrent())
	b = app(b, 'J', d.ExpLuminosity())
	b = app(b, 'N', d.ExpAmount())
	return b
}

// IsDimensionless returns true if d is dimensionless, that is to say all dimension exponents are zero.
func (d Dimension) IsDimensionless() bool { return d.dims == ([7]int16{}) }

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

type Prefix int8

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

// fixed point representation integer supported by this package.
type fixed interface {
	~int64
}

func formatAppend[T fixed](b []byte, value T, base Prefix, fmt byte, prec int) []byte {
	switch {
	case fmt != 'f':
		return append(b, "<si!INVALID FMT>"...)
	case prec < 0:
		return append(b, "<si!NEGATIVE PREC>"...)
	case value == 0:
		return append(b, '0')
	case !base.IsValid():
		return append(b, "<si!BAD BASE>"...)
	}

	isNegative := value < 0
	if isNegative {
		value = -value
	}
	v64 := int64(value)
	log10 := ilog10(v64)

	log10mod3 := log10 % 3
	frontDigits := (log10mod3) + 1
	allDigits := log10 + 1
	divDiff := log10 - log10mod3
	front := v64 / iLogTable[divDiff]
	if front == 0 {
		frontDigits = 0 // Underflowed format front.
	}
	var back int64
	if frontDigits == 0 || (allDigits > 3 && prec > frontDigits) {
		back = v64 % iLogTable[divDiff]
		backDigits := divDiff
		excess := backDigits + frontDigits - prec
		if excess > 0 {
			back /= iLogTable[excess]
		}
	}
	if isNegative {
		b = append(b, '-')
	}
	b = strconv.AppendInt(b, front, 10)
	if back != 0 {
		b = append(b, '.')
		b = strconv.AppendInt(b, back, 10)
	}
	// Calculate new base.
	base += Prefix(log10 - log10mod3)
	if base != PrefixNone {
		b = append(b, base.String()...)
	}
	return b
}

var iLogTable = [...]int64{
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

// ilog10 returns the integer logarithm base 10 of v, which
// can be interpreted as the quanity of digits in the number in base 10 minus one.
func ilog10(v int64) int {
	for i, l := range iLogTable {
		if v < l {
			return i - 1
		}
	}
	return len(iLogTable)
}
