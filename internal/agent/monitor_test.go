package agent

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricSet_Collect(t *testing.T) {
	m := &MetricSet{}

	for i := 0; i < 5; i += 1 {
		prevPollCount := m.PollCount
		prevRandomValue := m.RandomValue

		m.Collect()

		assert.Equal(t, prevPollCount+1, m.PollCount, fmt.Sprintf("MetricSet.Collect() error: MetricSet.PollCount = %v, want %v\n", m.PollCount, prevPollCount+1))

		assert.NotEqual(t, prevRandomValue, m.RandomValue, fmt.Sprintf("MetricSet.Collect error: MetricSet.RandomValue = %v, want different from before-collecting value\n", m.RandomValue))

		if m.Mallocs <= 0 {
			t.Errorf("MetricSet.Collect error: MetricSet.Mallocs = %v, want > 0 value\n", m.Mallocs)
		}
	}
}
