package srapi

import (
	"net/url"
	"strconv"
)

type Link struct {
	Relation string `json:"rel"`
	URI      string `json:"uri"`
}

func (self *Link) request() request {
	relURL := self.URI[len(BaseUrl):]

	return request{"GET", relURL, nil, nil}
}

type Pagination struct {
	Offset int
	Max    int
	Size   int
	Links  []Link
}

// for the 'hasLinks' interface
func (self *Pagination) links() []Link {
	return self.Links
}

type Cursor struct {
	Offset int
	Max    int
}

func (self *Cursor) applyToURL(u *url.URL) {
	values := u.Query()

	values.Set("offset", strconv.Itoa(self.Offset))
	values.Set("max", strconv.Itoa(self.Max))

	u.RawQuery = values.Encode()
}

type Direction int

const (
	Ascending Direction = iota
	Descending
)

type Sorting struct {
	OrderBy   string
	Direction Direction
}

func (self *Sorting) applyToURL(u *url.URL) {
	values := u.Query()
	dir := "asc"

	if self.Direction == Descending {
		dir = "desc"
	}

	values.Set("orderby", self.OrderBy)
	values.Set("direction", dir)

	u.RawQuery = values.Encode()
}