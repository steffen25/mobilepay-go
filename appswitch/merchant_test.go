package appswitch

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMobilePayTimestamp_UnmarshalJSON(t *testing.T) {
	{
		mpTime := MobilePayTimestamp{time.Now()}
		err := mpTime.UnmarshalJSON([]byte("invalid-timestamp-format"))

		assert.NotNil(t, err)
	}

	{
		mpTime := MobilePayTimestamp{time.Now()}
		err := mpTime.UnmarshalJSON([]byte("2016-04-08T07:45:36.533"))

		assert.Nil(t, err)
	}
}
