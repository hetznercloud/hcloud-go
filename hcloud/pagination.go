package hcloud

type pageGetter func(start, end int) (*Response, bool, error)

// Page serves as accessor of the API pagination.
type Page struct {
	pageGetter  pageGetter
	err         error
	currentPage int
	response    *Response
}

// Current returns the current page number.
func (p *Page) Current() int {
	return p.currentPage
}

// Response returns the Response of the API request.
func (p *Page) Response() *Response {
	return p.response
}

// GoTo fetches the given page.
// It returns true on success, or false if the resource is exhausted or an error happened while fetching it.
func (p *Page) GoTo(page int) bool {
	var exhausted bool
	p.response, exhausted, p.err = p.pageGetter(page, page)
	if p.err == nil {
		p.currentPage = page
	}
	return !exhausted
}

// Err returns the error, if any, that was encountered during iteration.
func (p *Page) Err() error {
	return p.err
}

// Next fetches the next page.
// It returns true on success, or false if the resource is exhausted or an error happened while fetching it.
func (p *Page) Next() bool {
	next := p.currentPage + 1
	var exhausted bool
	p.response, exhausted, p.err = p.pageGetter(next, next)
	if p.err == nil {
		p.currentPage = next
	}
	return !exhausted
}

// All fetches all pages.
func (p *Page) All() {
	p.response, _, p.err = p.pageGetter(0, 0)
	if p.err == nil {
		p.currentPage = 0
	}
}
