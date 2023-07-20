package internal

import (
	"errors"
	"fmt"
	"strconv"
)

// CheckDigit provides an interface for all check digit algorithms.
type CheckDigit interface {
	// Generate calculates a check digit
	Generate(seed string) (string, error)

	// Validate accepts the serial number and check digit portions of a tracking
	// number and validates the serial number.
	Validate(serial, code string) bool
}

// NewNoop instantiates a check digit algorithm that always passes.
func NewNoop() CheckDigit {
	return &Noop{}
}

func NewMod7() CheckDigit {
	return &Mod7{}
}

func NewMod10(em, om int) CheckDigit {
	return &Mod10{Multipliers: [2]int{em, om}}
}

func NewS10(weightings []int, exists []string) CheckDigit {
	if len(weightings) == 0 {
		weightings = []int{8, 6, 4, 2, 3, 5, 9, 7}
	}
	return &S10{
		Weightings: weightings,
	}
}

func NewMod3736() CheckDigit {
	return &Mod3736{}
}

func NewSumProductWithWeightingsAndModulo(weightings []int, m1, m2 int) CheckDigit {
	if m1 == 0 || m2 == 0 {
		return nil
	}
	return &SumProductWithWeightingsAndModulo{
		Weightings: weightings,
		Modulo1:    m1,
		Modulo2:    m2,
	}
}

// Noop is used if no check digit is needed.
type Noop struct{}

func (n Noop) Generate(s string) (string, error) {
	return "", nil
}

func (n Noop) Validate(s, cd string) bool {
	return true
}

type Mod7 struct{}

func (m Mod7) Generate(s string) (string, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return "0", nil
	}

	return fmt.Sprintf("%d", n%7), nil
}

func (m Mod7) Validate(s, cd string) bool {
	v, err := m.Generate(s)
	if err != nil {
		return false
	}

	return (v == cd)
}

type Mod10 struct {
	// 0: Even 1: Odd
	Multipliers [2]int
}

func (m10 *Mod10) Generate(s string) (string, error) {
	var sum int
	for i := 0; i < len(s); i++ {
		c := s[i]
		n, err := strconv.Atoi(string(c))
		if err != nil { // must be alpha
			n = int((c - 3) % 10)
		}
		sum += n * m10.Multipliers[i%2]
	}

	return fmt.Sprintf("%d", (10-sum%10)%10), nil
}

func (m10 *Mod10) Validate(s, cd string) bool {
	v, err := m10.Generate(s)
	if err != nil {
		return false
	}

	return (v == cd)
}

type S10 struct {
	Weightings []int
}

func (s10 *S10) Generate(s string) (string, error) {
	if len(s) != len(s10.Weightings) {
		return "", fmt.Errorf("seed should be same length as weightings")
	}

	var sum int
	for i, weight := range s10.Weightings {
		n, err := strconv.Atoi(string(s[i]))
		if err != nil {
			return "", err
		}
		sum += n * weight
	}

	sum = sum % 11
	switch sum {
	case 1:
		sum = 0
	case 0:
		sum = 5
	default:
		sum = 11 - sum
	}

	return fmt.Sprintf("%d", sum), nil
}

func (s10 *S10) Validate(s, cd string) bool {
	v, err := s10.Generate(s)
	if err != nil {
		return false
	}

	return (v == cd)
}

type SumProductWithWeightingsAndModulo struct {
	Weightings []int
	Modulo1    int
	Modulo2    int
}

func (sp *SumProductWithWeightingsAndModulo) Generate(s string) (string, error) {
	if len(s) < len(sp.Weightings) {
		return "", fmt.Errorf("seed should be at least as long as weightings, %d %d", len(s), len(sp.Weightings))
	}

	// only use the last N characters of the input string.
	s = s[len(s)-len(sp.Weightings):]

	var sum int
	for i, weight := range sp.Weightings {
		n, err := strconv.Atoi(string(s[i]))
		if err != nil {
			return "", err
		}
		sum += n * weight
	}

	return fmt.Sprintf("%d", sum%sp.Modulo1%sp.Modulo2), nil
}

func (sp *SumProductWithWeightingsAndModulo) Validate(s, cd string) bool {
	v, err := sp.Generate(s)
	if err != nil {
		return false
	}

	return (v == cd)
}

type Mod3736 struct{}

func (m3 *Mod3736) Generate(s string) (string, error) {
	const mod = 36
	var n byte = mod

	for i := 0; i < len(s); i++ {
		si, err := mod36CharToVal(s[i])
		if err != nil {
			continue
		}
		n += si
		if n > mod {
			n -= mod
		}
		n *= 2
		if n > mod {
			n = n - mod - 1
		}
	}

	n = mod + 1 - n
	n = n % 36

	return string(mod36ValtoChar(n)), nil
}

func (m3 *Mod3736) Validate(s, cd string) bool {
	v, err := m3.Generate(s)
	if err != nil {
		return false
	}

	return (v == cd)
}

func mod36CharToVal(b byte) (byte, error) {
	if b >= '0' && b <= '9' {
		return b - '0', nil
	}

	if b >= 'A' && b <= 'Z' {
		return b - 'A' + 10, nil
	}

	if b >= 'a' && b <= 'z' {
		return b - 'a' + 10, nil
	}

	return 255, errors.New("invalid character")
}

func mod36ValtoChar(b byte) byte {
	b %= 36
	if b < 10 {
		return b + '0'
	}

	return b - 10 + 'A'
}
