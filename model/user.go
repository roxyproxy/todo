package model

type User struct {
	Id        string `json:"id"`
	UserName  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"-"`
	Location  string `json:"location"` //time.Location
}

type Credentials struct {
	Password string `json:"password"`
	UserName string `json:"username"`
}
