package model

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type User struct {
	Id        string         `json:"id"`
	UserName  string         `json:"username"`
	FirstName string         `json:"firstname"`
	LastName  string         `json:"lastname"`
	Password  string         `json:"-"`
	Location  CustomLocation `json:"location"` // time.Location
}

type NewUser struct {
	UserName  string         `json:"username"`
	FirstName string         `json:"firstname"`
	LastName  string         `json:"lastname"`
	Password  string         `json:"password"`
	Location  CustomLocation `json:"location"` // time.Location
}

type Credentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	UserId string `json:"userid"`
	jwt.StandardClaims
}

type Token struct {
	TokenString string `json:"token"`
}

type KeyUserId string

type CustomLocation struct {
	*time.Location
}

func (c *CustomLocation) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		return nil
	}

	c.Location, err = time.LoadLocation(s)
	return
}

func (c CustomLocation) MarshalJSON() ([]byte, error) {
	if c.Location.String() == "" {
		return nil, nil
	}
	// return []byte(fmt.Sprintf(`"%s"`, c.Location.String())), nil
	return json.Marshal(c.Location.String())
}
