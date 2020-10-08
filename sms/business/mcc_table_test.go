package business

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountryCodeViaIMSI(t *testing.T) {
	unitCases := map[string]string{
		"":   "",
		"CN": "460092145800000",
		"HK": "454065400407000",
		"GB": "234336538092000",
		"DE": "262021915381000",
	}
	for expected, prefix := range unitCases {
		code, _ := CountryCodeViaIMSI(prefix)
		assert.Equal(t, expected, code)
	}
}

//noinspection SpellCheckingInspection
func TestFormatPhoneNumber(t *testing.T) {
	unitCases := [][]string{
		{"", "", ""},
		{"000000000000000", "1380013000", "1380013000"},
		{"460000000000000", "1380013000", "+861380013000"},
		{"460000000000000", "N/A", "N/A"},
	}
	for _, unitCase := range unitCases {
		var (
			imsi      = unitCase[0]
			input     = unitCase[1]
			expected  = unitCase[2]
			formatted = formatPhoneNumber(imsi, input)
		)
		assert.Equal(t, expected, formatted)
	}
}
