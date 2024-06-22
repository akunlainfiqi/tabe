package entities

type User struct {
	id         string
	givenName  string
	familyName string
	name       string
	email      string
	picture    string
	locale     string
}

func NewUser(id, givenName, familyName, name, email, picture, locale string) *User {
	return &User{
		id:         id,
		givenName:  givenName,
		familyName: familyName,
		name:       name,
		email:      email,
		picture:    picture,
		locale:     locale,
	}
}

func (u *User) ID() string {
	return u.id
}

func (u *User) GivenName() string {
	return u.givenName
}

func (u *User) FamilyName() string {
	return u.familyName
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Picture() string {
	return u.picture
}

func (u *User) Locale() string {
	return u.locale
}
