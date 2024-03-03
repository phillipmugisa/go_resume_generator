package storage

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/phillipmugisa/go_resume_generator/data"
)

type Storage interface {
	// user data
	CreateUser(data.User) error
	CreateUserimage(data.User, string) error
	GetUsers(map[string]string) ([]*data.User, error)
	DeleteUser(data.User) error
	VerifyUserEmail(string) ([]*data.User, error)
	SetUserSocials(u data.User) error

	// profile
	CreateProfile(data.Profile) error
	GetProfile(string) (*data.Profile, error)
	GetProfileByRole(string) ([]*data.Profile, error)
	DeleteProfile(data.Profile) error

	// Session
	CreateSession(data.Session) error
	GetSession(string) (*data.Session, error)
	DeleteSession(data.Session) error
	CancelSession(data.Session) error

	// Projects
	CreateProject(data.Project) error
	GetProjects(map[string]string) ([]*data.Project, error)
	GetProjectsByTechStack(map[string]string) ([]*data.Project, error)
	DeleteProject(int) error

	// Employment
	CreateEmployment(data.Employment) error
	GetEmployments(map[string]string) ([]*data.Employment, error)
	GetEmploymentsByTechStack(map[string]string) ([]*data.Employment, error)
	DeleteEmployment(int) error

	// Hobby
	CreateHobby(data.Hobby) error
	GetHobbies(map[string]string) ([]*data.Hobby, error)
	DeleteHobby(int) error

	// TechStack
	CreateTechStack(data.TechStack) error
	GetTechStacks(map[string]string) ([]*data.TechStack, error)
	GetProjectTechStacks(map[string]string) ([]*data.TechStack, error)
	GetEmploymentTechStacks(map[string]string) ([]*data.TechStack, error)
	AddTechStackToProject(data.TechStack, data.Project) error
	AddTechStackToEmployment(data.TechStack, data.Project) error
	DeleteTechStack(int) error
}

func scanUsers(rows *sql.Rows) ([]*data.User, error) {
	users := []*data.User{}
	for rows.Next() {
		user := new(data.User)
		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Firstname,
			&user.Lastname,
			&user.Email,
			&user.Bio,
			&user.Phone,
			&user.Country,
			&user.Password,
		)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GenerateRecordId() string {

	uuid := uuid.New()
	return uuid.String()
}
