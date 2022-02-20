package mobilepay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArgError_Error(t *testing.T) {
	err := newArgError("paymentId", "paymentId is empty")
	expected := "paymentId is invalid because paymentId is empty"
	assert.Error(t, err)
	assert.Equal(t, expected, err.Error())
}
