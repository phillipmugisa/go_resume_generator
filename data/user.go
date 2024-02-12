package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Bio       string `json:"bio"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Image     string `json:"image"`
	Country   string `json:"country"`

	// social
	Portfolio string `json:"portfolio"`
	Github    string `json:"github"`
	Linkedin  string `json:"linkedin"`
	Twitter   string `json:"twitter"`

	Start_date    time.Time `json:"start_date"` // period user started programming
	Years_of_work int       `json:"years_of_work"`

	Created_on     time.Time `json:"created_on"`
	Updated_on     time.Time `json:"updated_on"`
	Last_sign_in   time.Time `json:"last_sign_in"`
	Email_verified bool      `json:"email_verified"`
}

func (u User) String() string {
	return u.Username
}

func NewUser(firstname, lastname, username, email, password, phone, bio, country, start_date string) (*User, error) {
	pwd, err := HashPassword(password)
	if err != nil {
		return nil, errors.New("Error hashing password")
	}

	if start_date == "" {
		return nil, errors.New("Provide Programming Start Period.")
	}

	st, st_err := time.Parse("01/12/2024", start_date)
	if st_err != nil {
		return nil, st_err
	}

	years_of_work, err := GetWorkDuration(st, time.Now())
	if err != nil {
		return nil, err
	}

	return &User{
		Bio:            bio,
		Firstname:      firstname,
		Lastname:       lastname,
		Username:       username,
		Email:          email,
		Password:       pwd,
		Created_on:     time.Now(),
		Updated_on:     time.Now(),
		Start_date:     st,
		Years_of_work:  years_of_work,
		Email_verified: false,
		Phone:          phone,
		Country:        country,
	}, nil
}

func (u *User) SetSocials(portfolio, github, linkedin, twitter string) error {
	u.Portfolio = portfolio
	u.Github = github
	u.Linkedin = linkedin
	u.Twitter = twitter

	return nil
}

func HashPassword(p string) (string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pwd), nil
}

func (u *User) ValidateUser() error {
	// perform user data validations
	// this in called before savsing to database

	return nil
}

// user Profile
type Profile struct {
	User  User   `json:"user"`
	Role  string `json:"role"`
	About string `json:"about"`

	Views int `json:"views"`
}

func (u User) NewProfile(role, about string) (*Profile, error) {
	return &Profile{
		User:  u,
		Role:  role,
		About: about,
		Views: 0,
	}, nil
}
