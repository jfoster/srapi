package srapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Game struct {
	Id    string
	Names struct {
		International string
		Japanese      string
	}
	Abbreviation string
	Weblink      string
	Released     int
	Ruleset      struct {
		ShowMilliseconds    bool           `json:"show-milliseconds"`
		RequireVerification bool           `json:"require-verification"`
		RequireVideo        bool           `json:"require-video"`
		RunTimes            []TimingMethod `json:"run-times"`
		DefaultTime         TimingMethod   `json:"default-time"`
		EmulatorsAllowed    bool           `json:"emulators-allowed"`
	}
	Romhack bool

	// do not use this field directly, use the available methods
	PlatformsData interface{} `json:"platforms"`

	// do not use this field directly, use the available methods
	RegionsData interface{} `json:"regions"`

	// do not use this field directly, use the available methods
	ModeratorsData interface{} `json:"moderators"`

	// do not use this field directly, use the available methods
	CategoriesData interface{} `json:"categories"`

	// do not use this field directly, use the available methods
	LevelsData interface{} `json:"levels"`

	// do not use this field directly, use the available methods
	VariablesData interface{} `json:"variables"`

	Created string
	Assets  map[string]*AssetLink
	Links   []Link
}

type AssetLink struct {
	URI    string
	Width  int
	Height int
}

type gameResponse struct {
	Data Game
}

func GameById(id string) (*Game, *Error) {
	request := request{"GET", "/games/" + id, nil, nil, nil}
	result := &gameResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func GameByAbbreviation(abbrev string) (*Game, *Error) {
	return GameById(abbrev)
}

func (self *Game) PlatformIds() []string {
	result := make([]string, 0)

	switch asserted := self.PlatformsData.(type) {
	// list of IDs (strings)
	case []interface{}:
		for _, something := range asserted {
			id, okay := something.(string)
			if okay {
				result = append(result, id)
			}
		}

	// sub-resource due to embeds, aka "{data:....}"
	// TODO: skip the conversion back and forth and just assert our way through the available data
	case map[string]interface{}:
		for _, platform := range self.Platforms() {
			result = append(result, platform.Id)
		}
	}

	return result
}

func (self *Game) Platforms() []*Platform {
	result := make([]*Platform, 0)

	switch asserted := self.PlatformsData.(type) {
	// list of IDs (strings)
	case []interface{}:
		for _, id := range self.PlatformIds() {
			platform, err := PlatformById(id)
			if err == nil {
				result = append(result, platform)
			}
		}

	// sub-resource due to embeds, aka "{data:....}"
	case map[string]interface{}:
		// convert generic mess into JSON
		encoded, _ := json.Marshal(asserted)

		// ... and try to turn it back into something meaningful
		dest := PlatformCollection{}
		err := json.Unmarshal(encoded, &dest)
		if err == nil {
			for idx := range dest.Data {
				result = append(result, &dest.Data[idx])
			}
		}
	}

	return result
}

func (self *Game) RegionIds() []string {
	result := make([]string, 0)

	switch asserted := self.RegionsData.(type) {
	// list of IDs (strings)
	case []interface{}:
		for _, something := range asserted {
			id, okay := something.(string)
			if okay {
				result = append(result, id)
			}
		}

	// sub-resource due to embeds, aka "{data:....}"
	// TODO: skip the conversion back and forth and just assert our way through the available data
	case map[string]interface{}:
		for _, region := range self.Regions() {
			result = append(result, region.Id)
		}
	}

	return result
}

func (self *Game) Regions() []*Region {
	result := make([]*Region, 0)

	switch asserted := self.RegionsData.(type) {
	// list of IDs (strings)
	case []interface{}:
		for _, id := range self.RegionIds() {
			region, err := RegionById(id)
			if err == nil {
				result = append(result, region)
			}
		}

	// sub-resource due to embeds, aka "{data:....}"
	case map[string]interface{}:
		// convert generic mess into JSON
		encoded, _ := json.Marshal(asserted)

		// ... and try to turn it back into something meaningful
		dest := RegionCollection{}
		err := json.Unmarshal(encoded, &dest)
		if err == nil {
			for idx := range dest.Data {
				result = append(result, &dest.Data[idx])
			}
		}
	}

	return result
}

// for the 'hasLinks' interface
func (self *Game) links() []Link {
	return self.Links
}

type GameCollection struct {
	Data       []Game
	Pagination Pagination
}

type GameFilter struct {
	Name         string
	Abbreviation string
	Released     int
	Platform     string
	Region       string
	Moderator    string
	Romhack      *bool
}

func (self *GameFilter) applyToURL(u *url.URL) {
	values := u.Query()

	if len(self.Name) > 0 {
		values.Set("name", self.Name)
	}

	if len(self.Abbreviation) > 0 {
		values.Set("abbreviation", self.Abbreviation)
	}

	if self.Released > 0 {
		values.Set("released", strconv.Itoa(self.Released))
	}

	if len(self.Platform) > 0 {
		values.Set("platform", self.Platform)
	}

	if len(self.Region) > 0 {
		values.Set("region", self.Region)
	}

	if len(self.Moderator) > 0 {
		values.Set("moderator", self.Moderator)
	}

	if self.Romhack != nil {
		if *self.Romhack {
			values.Set("romhack", "yes")
		} else {
			values.Set("romhack", "no")
		}
	}

	u.RawQuery = values.Encode()
}

func Games(f *GameFilter, s *Sorting, c *Cursor) (*GameCollection, *Error) {
	return fetchGames(request{"GET", "/games", f, s, c})
}

func (self *GameCollection) NextPage() (*GameCollection, *Error) {
	return self.fetchLink("next")
}

func (self *GameCollection) PrevPage() (*GameCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *GameCollection) fetchLink(name string) (*GameCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchGames(next.request())
}

func fetchGames(request request) (*GameCollection, *Error) {
	result := &GameCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
