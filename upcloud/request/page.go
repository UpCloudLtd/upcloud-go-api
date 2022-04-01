package request

import (
	"fmt"
	"net/url"
)

const (
	PageSizeMax       uint16 = 100
	PageResultMaxSize int    = 1000
)

var DefaultPage = &Page{Number: 1, Size: PageSizeMax}

type Page struct {
	Size   uint16
	Number uint16
}

func (p *Page) offset() uint16 {
	if p.Number < 2 {
		return 0
	}
	return (p.Number - 1) * p.Size
}

func (p *Page) SizeInt() int {
	return int(p.Size)
}

func (p *Page) NumberInt() int {
	return int(p.Number)
}

func (p *Page) Values() url.Values {
	if p.Number < 1 {
		p.Number = 1
	}
	v := url.Values{}
	v.Add("limit", fmt.Sprint(p.Size))
	v.Add("offset", fmt.Sprint(p.offset()))
	return v
}

func (p *Page) Next() *Page {
	return &Page{
		Size:   p.Size,
		Number: p.Number + 1,
	}
}

func (p *Page) Previous() *Page {
	var n uint16
	if p.Number > 1 {
		n = p.Number - 1
	}

	return &Page{
		Size:   p.Size,
		Number: n,
	}
}

func (p *Page) String() string {
	return p.Values().Encode()
}
