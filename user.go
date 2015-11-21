package srapi

import "net/url"

type User struct {
	Id    string
	Names struct {
		International string
		Japanese      string
	}
	Weblink   string
	NameStyle struct {
		Style     string
		Color     *NameColor
		ColorFrom *NameColor `json:"color-from"`
		ColorTo   *NameColor `json:"color-to"`
	} `json:"name-style"`
	Role     string
	Signup   string
	Location struct {
		Country Location
		Region  *Location
	}
	Twitch        *SocialLink
	Hitbox        *SocialLink
	YouTube       *SocialLink
	Twitter       *SocialLink
	SpeedRunsLive *SocialLink
	Links         []Link
}

type SocialLink struct {
	URI string
}

type Location struct {
	Code  string
	Names struct {
		International string
		Japanese      string
	}
}

type NameColor struct {
	Light string
	Dark  string
}

type userResponse struct {
	Data User
}

func UserById(id string) (*User, *Error) {
	return fetchUser(request{"GET", "/users/" + id, nil, nil, nil})
}

func fetchUser(request request) (*User, *Error) {
	result := &userResponse{}

	err := httpClient.do(request, result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// for the 'hasLinks' interface
func (self *User) links() []Link {
	return self.Links
}

type UserCollection struct {
	Data       []User
	Pagination Pagination
}

func (self *UserCollection) users() []*User {
	result := make([]*User, 0)

	for idx := range self.Data {
		result = append(result, &self.Data[idx])
	}

	return result
}

type UserFilter struct {
	Lookup        string
	Name          string
	Twitch        string
	Hitbox        string
	Twitter       string
	SpeedRunsLive string
}

func (self *UserFilter) applyToURL(u *url.URL) {
	values := u.Query()

	if len(self.Lookup) > 0 {
		values.Set("lookup", self.Lookup)
	}

	if len(self.Name) > 0 {
		values.Set("name", self.Name)
	}

	if len(self.Twitch) > 0 {
		values.Set("twitch", self.Twitch)
	}

	if len(self.Hitbox) > 0 {
		values.Set("hitbox", self.Hitbox)
	}

	if len(self.Twitter) > 0 {
		values.Set("twitter", self.Twitter)
	}

	if len(self.SpeedRunsLive) > 0 {
		values.Set("speedrunslive", self.SpeedRunsLive)
	}

	u.RawQuery = values.Encode()
}

func Users(f *UserFilter, s *Sorting, c *Cursor) (*UserCollection, *Error) {
	return fetchUsers(request{"GET", "/users", f, s, c})
}

func (self *UserCollection) NextPage() (*UserCollection, *Error) {
	return self.fetchLink("next")
}

func (self *UserCollection) PrevPage() (*UserCollection, *Error) {
	return self.fetchLink("prev")
}

func (self *UserCollection) fetchLink(name string) (*UserCollection, *Error) {
	next := firstLink(&self.Pagination, name)
	if next == nil {
		return nil, nil
	}

	return fetchUsers(next.request())
}

// always returns a collection, even when an error is returned;
// makes other code more monadic
func fetchUsers(request request) (*UserCollection, *Error) {
	result := &UserCollection{}

	err := httpClient.do(request, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
