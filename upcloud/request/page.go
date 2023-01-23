package request

import (
	"fmt"
	"net/url"
)

const (
	PageSizeMax       int = 100
	PageResultMaxSize int = 1000
)

var DefaultPage = &Page{Number: 1, Size: PageSizeMax}

// Page represents list query filter that can be used to limit the number of results that are returned from the endpoint.
type Page struct {
	// Size of the page limits the number of items in the list
	Size int
	// Page number sets offset where to start returning data.
	Number int
}

func (p *Page) offset() int {
	if p.Number < 1 {
		p.Number = 1
	}
	return (p.Number - 1) * p.Size
}

func (p *Page) limit() int {
	if p.Size < 1 {
		p.Size = 0
	}
	return p.Size
}

func (p *Page) Values() url.Values {
	v := url.Values{}
	v.Add("limit", fmt.Sprint(p.limit()))
	v.Add("offset", fmt.Sprint(p.offset()))
	return v
}

func (p *Page) Next() *Page {
	return &Page{
		Size:   p.limit(),
		Number: p.Number + 1,
	}
}

func (p *Page) Previous() *Page {
	var n int
	if p.Number > 1 {
		n = p.Number - 1
	} else {
		n = 1
	}

	return &Page{
		Size:   p.Size,
		Number: n,
	}
}

func (p *Page) String() string {
	return p.Values().Encode()
}

func (p *Page) ToQueryParam() string {
	return p.String()
}
