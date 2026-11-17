package money

type Paise int64

const (
	RupeePaise = 100
	BasisPoint = 10000
)

func (p Paise) Int64() int64 {
	return int64(p)
}

func (p Paise) String() string {
	rupees := p / RupeePaise
	paise := p % RupeePaise
	if paise < 0 {
		paise = -paise
	}
	return fmt.Sprintf("%d.%02d", rupees, paise)
}

func (p Paise) MulBP(bp int64) Paise {
	return Paise((p.Int64()*bp + 5000) / 10000)
}

func (p Paise) MulFrac(num, den int64) Paise {
	if den == 0 {
		return 0
	}
	return Paise((p.Int64()*num*2 + den) / (den * 2))
}

func (p Paise) RoundToRupee() Paise {
	return Paise((p.Int64() + 50) / 100 * 100)
}

func (p Paise) CeilToRupee() Paise {
	return Paise((p.Int64() + 99) / 100 * 100)
}

func (p Paise) Min(other Paise) Paise {
	if p < other {
		return p
	}
	return other
}

func (p Paise) Max(other Paise) Paise {
	if p > other {
		return p
	}
	return other
}

func PaiseFromRupees(rupees float64) Paise {
	return Paise(int64(rupees * RupeePaise + 0.5))
}

func PaiseFromString(s string) (Paise, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty string")
	}

	negative := false
	if s[0] == '-' {
		negative = true
		s = s[1:]
	}

	parts := strings.Split(s, ".")
	rupees := int64(0)
	paise := int64(0)

	if len(parts) > 0 && parts[0] != "" {
		var err error
		rupees, err = strconv.ParseInt(strings.ReplaceAll(parts[0], ",", ""), 10, 64)
		if err != nil {
			return 0, err
		}
	}

	if len(parts) > 1 {
		frac := parts[1]
		if len(frac) > 2 {
			frac = frac[:2]
		}
		for len(frac) < 2 {
			frac += "0"
		}
		var err error
		paise, err = strconv.ParseInt(frac, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	total := rupees*RupeePaise + paise
	if negative {
		total = -total
	}
	return Paise(total), nil
}

func (p Paise) FormatINR() string {
	return p.String()
}