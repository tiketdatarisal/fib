package ursa

type User struct {
	Email        string  `json:"email"`
	Name         string  `json:"name,omitempty"`
	GoogleToken  string  `json:"googleToken,omitempty"`
	Organization string  `json:"organization,omitempty"`
	Groups       []Group `json:"groups,omitempty"`
	Deactivated  bool    `json:"deactivated"`
}

// NewUser create a new user using supplied values.
func NewUser(email, name, googleToken, organization, status string, group []Group) User {
	return User{
		Email:        email,
		Name:         name,
		GoogleToken:  googleToken,
		Organization: organization,
		Groups:       group,
		Deactivated:  func() bool { return status != "active" }(),
	}
}
