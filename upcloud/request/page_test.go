package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	p := &Page{
		Size:   100,
		Number: 0,
	}
	assert.Equal(t, "limit=100&offset=0", p.String())
	assert.Equal(t, uint16(1), p.Number)
	p = p.Next()
	assert.Equal(t, "limit=100&offset=100", p.String())
	assert.Equal(t, uint16(2), p.Number)
	p = p.Next()
	assert.Equal(t, "limit=100&offset=200", p.String())
	assert.Equal(t, uint16(3), p.Number)
	p = p.Next()
	assert.Equal(t, "limit=100&offset=300", p.String())
	assert.Equal(t, uint16(4), p.Number)
	p = p.Previous()
	assert.Equal(t, "limit=100&offset=200", p.String())
	assert.Equal(t, uint16(3), p.Number)
	p = &Page{
		Size:   100,
		Number: 1,
	}
	assert.Equal(t, "limit=100&offset=0", p.String())
	p = p.Previous()
	assert.Equal(t, "limit=100&offset=0", p.String())
}
