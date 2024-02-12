package data

import (
	"time"
)

// represents programming languages and tools
type TechStack struct {
	User User   `json:"user"`
	Name string `json:"name"`
}

func (u User) NewTechStack(n string) *TechStack {
	return &TechStack{
		User: u,
		Name: n,
	}
}

func (s TechStack) String() string {
	return s.Name
}

type Project struct {
	User        User        `json:"user"`
	Name        string      `json:"name"`
	Duration    int         `json:"duration"`
	Start_date  time.Time   `json:"start_date"`
	End_date    time.Time   `json:"end_date"`
	Status      string      `json:"status"`
	Github      string      `json:"github"`
	Prod_link   string      `json:"prod_link"`
	Description string      `json:"description"`
	Stack       []TechStack `json:"stack"`

	Created_on time.Time `json:"created_on"`
	Updated_on time.Time `json:"updated_on"`
}

func (u User) NewProject(name, status, github, prod_link, description string, start, end string) (*Project, error) {

	st, st_err := time.Parse("01/12/2024", start)
	if st_err != nil {
		return nil, st_err
	}
	et, et_err := time.Parse("01/12/2024", end)
	if et_err != nil {
		return nil, et_err
	}

	duration, err := GetWorkDuration(st, et)
	if err != nil {
		return nil, err
	}

	return &Project{
		User:        u,
		Name:        name,
		Duration:    duration,
		Status:      status,
		Start_date:  st,
		End_date:    et,
		Stack:       []TechStack{},
		Github:      github,
		Prod_link:   prod_link,
		Description: description,
		Created_on:  time.Now(),
		Updated_on:  time.Now(),
	}, nil
}

func (p Project) String() string {
	return p.Name
}

func (p *Project) AddStack(s TechStack) error {
	p.Stack = append(p.Stack, s)
	return nil
}

type Employment struct {
	User        User        `json:"user"`
	Name        string      `json:"name"`
	Employee    string      `json:"employee"`
	Start_date  time.Time   `json:"start_date"`
	End_date    time.Time   `json:"end_date"`
	Status      string      `json:"status"`
	Stack       []TechStack `json:"stack"`
	Prod_link   string      `json:"prod_link"`
	Duration    int         `json:"duration"`
	Description string      `json:"description"`

	Created_on time.Time `json:"created_on"`
	Updated_on time.Time `json:"updated_on"`
}

func (e Employment) String() string {
	return e.Name
}

func (u User) NewEmployment(name, status, prod_link, description string, start, end string) (*Employment, error) {

	st, st_err := time.Parse("01/12/2024", start)
	if st_err != nil {
		return nil, st_err
	}
	et, et_err := time.Parse("01/12/2024", end)
	if et_err != nil {
		return nil, et_err
	}

	duration, err := GetWorkDuration(st, et)
	if err != nil {
		return nil, err
	}

	return &Employment{
		User:        u,
		Name:        name,
		Duration:    duration,
		Status:      status,
		Start_date:  st,
		End_date:    et,
		Stack:       []TechStack{},
		Prod_link:   prod_link,
		Description: description,
		Created_on:  time.Now(),
		Updated_on:  time.Now(),
	}, nil
}

func (e *Employment) AddStack(s TechStack) error {
	e.Stack = append(e.Stack, s)
	return nil
}

type Hobby struct {
	User User   `json:"user"`
	Name string `json:"name"`
}

func (u User) NewHobby(n string) *Hobby {
	return &Hobby{
		User: u,
		Name: n,
	}
}

func (h Hobby) String() string {
	return h.Name
}
