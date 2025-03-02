package utilsValidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateNanoid(t *testing.T) {
	// Test a valid nanoid.
	validNanoid := "V1StGXR8_Z5jdHi6B-myT"
	assert.True(t, ValidateNanoid(validNanoid))

	// Test an invalid nanoid.
	invalidNanoid := "invalidnanoid"
	assert.False(t, ValidateNanoid(invalidNanoid))
}

func TestValidateNanoidArray(t *testing.T) {
	// Test a valid nanoid array.
	validNanoidArray := []string{"V1StGXR8_Z5jdHi6B-myT", "V1StGXR8_Z5jdHi6B-myT"}
	assert.True(t, ValidateNanoidArray(validNanoidArray))

	// Test an invalid nanoid array.
	invalidNanoidArray := []string{"V1StGXR8_Z5jdHi6B-myT", "invalidnanoid"}
	assert.False(t, ValidateNanoidArray(invalidNanoidArray))
}
