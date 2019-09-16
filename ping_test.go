package ping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	_, err := Run("asdbasdfasdfdasfkjhasdfkjh.com", nil)
	assert.Error(t, err)

	stats, err := Run("www.baidu.com", &Opts{
		Count: 10,
	})
	if assert.NoError(t, err) {
		assert.True(t, stats.RTTMin > 0)
		assert.True(t, stats.RTTAvg > 0)
		assert.True(t, stats.RTTMax > 0)
		assert.True(t, stats.RTTDev > 0)
	}
}
