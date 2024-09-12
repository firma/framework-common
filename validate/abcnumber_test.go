package validate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateAAbc(t *testing.T) {
	checks := checkAbc("bdfdsfdsfsdfs")
	assert.Equal(t, checks, true)

	checks2 := checkAbc("bdfdsf1dsfsdfs")
	assert.Equal(t, checks2, false)
}

func TestValidateAbcNumber(t *testing.T) {
	checks := checkNumber("12312312")
	assert.Equal(t, checks, true)
	checks2 := checkNumber("1s2312312")
	assert.Equal(t, checks2, false)
}

func TestValidateTell(t *testing.T) {
	checks := checkTell("1121-312312")
	assert.Equal(t, checks, true)
	checks2 := checkTell("111s21-312312")
	assert.Equal(t, checks2, false)
}
