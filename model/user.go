package model

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// User represents users structure.
type User struct {
	ID        string         `json:"id"`
	UserName  string         `json:"username"`
	FirstName string         `json:"firstname"`
	LastName  string         `json:"lastname"`
	Password  string         `json:"-"`
	Location  CustomLocation `json:"location"` // time.Location
}

// NewUser represents users structure with password.
type NewUser struct {
	UserName  string         `json:"username"`
	FirstName string         `json:"firstname"`
	LastName  string         `json:"lastname"`
	Password  string         `json:"password"`
	Location  CustomLocation `json:"location"` // time.Location
}

// Credentials represents users credentials.
type Credentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// Claims represents users claims in token string.
type Claims struct {
	UserID string `json:"userid"`
	jwt.StandardClaims
}

// Token represents token string.
type Token struct {
	TokenString string `json:"token"`
}

// KeyUserID represents userid in context.
type KeyUserID string

// CustomLocation used for Location json encodeing/decoding.
type CustomLocation struct {
	*time.Location
}

// UnmarshalJSON used to marshal/unmarshal CustomLocation.
func (c *CustomLocation) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		return nil
	}

	c.Location, err = time.LoadLocation(s)
	return
}

// MarshalJSON used to marshal/unmarshal CustomLocation.
func (c CustomLocation) MarshalJSON() ([]byte, error) {
	if c.Location.String() == "" {
		return nil, nil
	}
	// return []byte(fmt.Sprintf(`"%s"`, c.Location.String())), nil
	return json.Marshal(c.Location.String())
}
