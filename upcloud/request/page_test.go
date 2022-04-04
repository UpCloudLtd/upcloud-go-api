package request

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	p := &Page{
		Size:   100,
		Number: 0,
	}
	assert.Equal(t, "limit=100&offset=0", p.String())
	assert.Equal(t, 1, p.Number)
	p = p.Next()
	assert.Equal(t, "limit=100&offset=100", p.String())
	assert.Equal(t, 2, p.Number)
	p = p.Next()
	assert.Equal(t, "limit=100&offset=200", p.String())
	assert.Equal(t, 3, p.Number)
	p = p.Next()
	assert.Equal(t, "limit=100&offset=300", p.String())
	assert.Equal(t, 4, p.Number)
	p = p.Previous()
	assert.Equal(t, "limit=100&offset=200", p.String())
	assert.Equal(t, 3, p.Number)
	p = &Page{
		Size:   100,
		Number: 1,
	}
	assert.Equal(t, "limit=100&offset=0", p.String())
	p = p.Previous()
	assert.Equal(t, "limit=100&offset=0", p.String())
	p = &Page{
		Size:   100,
		Number: math.MaxInt64, // update to math.MaxInt when available (go 1.17)
	}
	p = p.Next()
	assert.Equal(t, "limit=100&offset=0", p.String())
	assert.Equal(t, 1, p.Number)
	p = p.Next()
	assert.Equal(t, "limit=100&offset=100", p.String())
	assert.Equal(t, 2, p.Number)
}
