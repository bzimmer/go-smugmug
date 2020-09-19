package smugmug

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Pages(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	p := &Pages{
		Total:          435,
		Start:          8,
		Count:          5,
		RequestedCount: 5,
	}

	a.Equal(8+5, p.Next())
	a.Equal(8-5, p.Previous())
	a.Equal(435-8-5, p.Remaining())
}
