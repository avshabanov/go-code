package logic

import (
	"fmt"
	"io"
	"time"
)

// Roles defines allowable values for user profile roles
var Roles = [...]string{
	"ADMIN",
	"EDITOR",
	"MODERATOR",
	"READER",
}

// UserProfile represents user account
type UserProfile struct {
	ID       int
	Name     string
	Created  time.Time
	Roles    []string
	Accounts []*OauthAccount
}

func (p *UserProfile) String() string {
	return fmt.Sprintf(
		"{id: %d, name: '%s', created: '%s', roles: %s, accounts: %s}",
		p.ID,
		p.Name,
		p.Created,
		p.Roles,
		p.Accounts,
	)
}

// OauthAccount represents user's oauth account
type OauthAccount struct {
	Token    string
	Provider string
	Created  time.Time
}

func (p *OauthAccount) String() string {
	return fmt.Sprintf(
		"{token: '%s', provider: '%s', created: '%s'}",
		p.Token,
		p.Provider,
		p.Created,
	)
}

// UserPage represents paginated result, that contains a collection of user profiles
type UserPage struct {
	Profiles    []*UserProfile
	OffsetToken string
}

// Dao represents an interface to user DAO
type Dao interface {
	io.Closer

	Add(profiles []*UserProfile) error
	QueryUsers(offsetToken string, limit int) (*UserPage, error)
	Get(id int) (*UserProfile, error)
	GetIDRange() (from int, to int, err error)
}
