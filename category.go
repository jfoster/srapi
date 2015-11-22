// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import "net/url"

// Category is a structure representing a game category, either full-game or per-level.
type Category struct {
	Id      string
	Name    string
	Weblink string
	Type    string
	Rules   string
	Players struct {
		Type  string
		Value int
	}
	Miscellaneous bool
	Links         []Link

	// do not use this field directly, use the available methods
	GameData interface{} `json:"game"`

	// do not use this field directly, use the available methods
	VariablesData interface{} `json:"variables"`
}

// toCategory transforms a data blob to a Category struct, if possible.
// Returns nil if casting the data was not successful or if data was nil.
func toCategory(data interface{}) *Category {
	dest := Category{}

	if data != nil && recast(data, &dest) == nil {
		return &dest
	}

	return nil
}

// toCategoryCollection transforms a data blob to a CategoryCollection.
// If data is nil or casting was unsuccessful, an empty CategoryCollection
// is returned.
func toCategoryCollection(data interface{}) *CategoryCollection {
	tmp := &CategoryCollection{}
	recast(data, tmp)

	return tmp
}

type categoryResponse struct {
	Data Category
}

// CategoryById tries to fetch a single category, identified by its Id.
// When an error is returned, the returned category is nil.
func CategoryById(id string) (*Category, *Error) {
	return fetchCategory(request{"GET", "/categories/" + id, nil, nil, nil})
}

// Game extracts the embedded game, if possible, otherwise it will fetch the
// game by doing one additional request. If nothing on the server side is fubar,
// then this function should never return nil.
func (c *Category) Game() *Game {
	if c.GameData == nil {
		return fetchGameLink(firstLink(c, "game"))
	}

	return toGame(c.GameData)
}

// Variables extracts the embedded variables, if possible, otherwise it will
// fetch them by doing one additional request. sort is only relevant when the
// variables are not already embedded.
func (c *Category) Variables(sort *Sorting) []*Variable {
	var collection *VariableCollection

	if c.VariablesData == nil {
		collection = fetchVariablesLink(firstLink(c, "variables"), nil, sort)
	} else {
		collection = toVariableCollection(c.VariablesData)
	}

	return collection.variables()
}

// PrimaryLeaderboard fetches the primary leaderboard, if any, for the category.
// The result can be nil.
func (c *Category) PrimaryLeaderboard(options *LeaderboardOptions) *Leaderboard {
	return fetchLeaderboardLink(firstLink(c, "leaderboard"), options)
}

// Records fetches a list of leaderboards for the category. For full-game
// categories, the list will contain one leaderboard, otherwise it will have one
// per level. This function always returns a LeaderboardCollection.
func (c *Category) Records(filter *LeaderboardFilter) *LeaderboardCollection {
	return fetchLeaderboardsLink(firstLink(c, "records"), filter, nil)
}

// Runs fetches a list of runs done in the given category, optionally filtered
// and sorted. This function always returns a RunCollection.
func (c *Category) Runs(filter *RunFilter, sort *Sorting) *RunCollection {
	return fetchRunsLink(firstLink(c, "records"), filter, sort)
}

// for the 'hasLinks' interface
func (c *Category) links() []Link {
	return c.Links
}

// CategoryCollection is one page of the entire category list. It consists of the
// categories as well as some pagination information (like links to the next or
// previous page).
type CategoryCollection struct {
	Data       []Category
	Pagination Pagination
}

// categories returns a list of pointers to the categories; used for cases where
// there is no pagination and the caller wants to return a flat slice of categories
// instead of a collection (which would be misleading, as collections imply
// pagination).
func (cc *CategoryCollection) categories() []*Category {
	var result []*Category

	for idx := range cc.Data {
		result = append(result, &cc.Data[idx])
	}

	return result
}

// CategoryFilter represents the possible filtering options when fetching a list
// of categories.
type CategoryFilter struct {
	Miscellaneous *bool
}

// applyToURL merged the filter into a URL.
func (cf *CategoryFilter) applyToURL(u *url.URL) {
	values := u.Query()

	if cf.Miscellaneous != nil {
		if *cf.Miscellaneous {
			values.Set("miscellaneous", "yes")
		} else {
			values.Set("miscellaneous", "no")
		}
	}

	u.RawQuery = values.Encode()
}

// NextPage tries to follow the "next" link and retrieve the next page of
// categories. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (cc *CategoryCollection) NextPage() (*CategoryCollection, *Error) {
	return cc.fetchLink("next")
}

// PrevPage tries to follow the "prev" link and retrieve the previous page of
// categories. If there is no such link, an empty collection and an error
// is returned. Otherwise, the error is nil.
func (cc *CategoryCollection) PrevPage() (*CategoryCollection, *Error) {
	return cc.fetchLink("prev")
}

// fetchLink tries to fetch a link, if it exists. If there is no such link, an
// empty collection and an error is returned. Otherwise, the error is nil.
func (cc *CategoryCollection) fetchLink(name string) (*CategoryCollection, *Error) {
	next := firstLink(&cc.Pagination, name)
	if next == nil {
		return &CategoryCollection{}, &Error{"", "", ErrorNoSuchLink, "Could not find a '" + name + "' link."}
	}

	return fetchCategories(next.request(nil, nil))
}

// fetchCategory fetches a single category from the network. If the request failed,
// the returned category is nil. Otherwise, the error is nil.
func fetchCategory(request request) (*Category, *Error) {
	result := &categoryResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// fetchCategoryLink tries to fetch a given link and interpret the response as
// a single category. If the link is nil or the category could not be fetched,
// nil is returned.
func fetchCategoryLink(link *Link) *Category {
	if link == nil {
		return nil
	}

	category, _ := fetchCategory(link.request(nil, nil))
	return category
}

// fetchCategories fetches a list of categories from the network. It always
// returns a collection, even when an error is returned.
func fetchCategories(request request) (*CategoryCollection, *Error) {
	result := &CategoryCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// fetchCategoriesLink tries to fetch a given link and interpret the response as
// a list of categories. It always returns a collection, even when an error is
// returned or the given link is nil.
func fetchCategoriesLink(link *Link, filter filter, sort *Sorting) *CategoryCollection {
	if link == nil {
		return &CategoryCollection{}
	}

	collection, _ := fetchCategories(link.request(filter, sort))
	return collection
}
