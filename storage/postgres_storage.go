package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/phillipmugisa/go_resume_generator/data"
)

type Storage interface {
	// user data
	CreateUser(data.User) error
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
	CreateSession(session data.Session) error
	GetSession(data.User) (*data.Session, error)
	DeleteSession(data.Session) error
	CancelSession(data.User) error
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}
	fmt.Println("Database Connection Successful...")
	return &PostgresStorage{
		db: db,
	}, nil
}

func initDB() (*sql.DB, error) {
	HOST := os.Getenv("POSTGRES_HOST")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DB")
	PORT := os.Getenv("POSTGRES_PORT")
	username := os.Getenv("POSTGRES_USER")

	// make db connection
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, errors.New("Couldnot connect to database")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// create required tables in db: Users
func (s *PostgresStorage) SetUpDB() error {
	if err := s.createUserTable(); err != nil {
		return err
	}

	if err := s.createProfileTable(); err != nil {
		return err
	}

	if err := s.createSessionTable(); err != nil {
		return err
	}

	return nil
}

// Users

func (s *PostgresStorage) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS Users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		firstname VARCHAR(255),
		lastname VARCHAR(255),
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		bio TEXT NULL,
		phone VARCHAR(20) NULL,
		country VARCHAR(100) NOT NULL,
		email_verified BOOLEAN DEFAULT FALSE,
		created_on TIMESTAMP NOT NUL,
		updated_on TIMESTAMP NOT NUL,
		last_sign_in TIMESTAMP NUL,
		portfolio VARCHAR(255) NULL,
		github VARCHAR(255) NULL,
		linkedin VARCHAR(255) NULL,
		twitter VARCHAR(255) NULL,
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) GetUsers(keywords map[string]string) ([]*data.User, error) {
	username, ok := keywords["username"]
	if !ok {
		username = ""
	}

	email, ok := keywords["email"]
	if !ok {
		email = ""
	}

	query := `SELECT username, firstname, lastname, email, bio, phone, country FROM Users WHERE username = $1 OR email = $2`
	rows, err := s.db.Query(query, username, email)
	if err != nil {
		return nil, err
	}
	return scanUsers(rows)
}
func (s *PostgresStorage) CreateUser(u data.User) error {

	query := `INSERT INTO Users (username, firstname, lastname, email, password, phone, country, created_on, updated_on, bio)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	_, err := s.db.Query(
		query,
		u.Username,
		u.Firstname,
		u.Lastname,
		u.Email,
		u.Password,
		u.Phone,
		u.Country,
		u.Created_on,
		u.Updated_on,
		u.Bio,
	)
	return err
}

func (s *PostgresStorage) SetUserSocials(u data.User) error {

	user_id, f_err := s.getUserID(u.Username)
	if f_err != nil {
		return nil
	}

	q := "UPDATE Users SET portfolio = $1, github = $2, linkedin = $3, twitter = $4 WHERE id = $5"

	_, err := s.db.Exec(q, u.Portfolio, u.Github, u.Linkedin, u.Twitter, user_id)
	return err
}

func (s *PostgresStorage) VerifyUserEmail(username string) ([]*data.User, error) {
	query := `UPDATE Users SET email_verified = TRUE  WHERE username = $1`
	_, err := s.db.Exec(query, username)
	if err != nil {
		return nil, err
	}
	return s.GetUsers(map[string]string{"username": username})
}

func (s *PostgresStorage) DeleteUser(u data.User) error {
	query := `DELETE FROM Users WHERE username = $1`
	_, err := s.db.Exec(query, u.Username)
	return err
}

func scanUsers(rows *sql.Rows) ([]*data.User, error) {
	users := []*data.User{}
	for rows.Next() {
		user := new(data.User)
		err := rows.Scan(
			&user.Username,
			&user.Firstname,
			&user.Lastname,
			&user.Email,
			&user.Bio,
			&user.Phone,
			&user.Country,
		)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *PostgresStorage) getUserID(username string) (int, error) {
	var user_id int
	f_err := s.db.QueryRow("SELECT id FROM Users WHERE username = $1", username).Scan(&user_id)
	if f_err != nil {
		return 0, f_err
	}
	return user_id, nil
}

// Users

// Profile

func (s *PostgresStorage) createProfileTable() error {
	query := `CREATE TABLE IF NOT EXISTS Profiles (
		id SERIAL PRIMARY KEY,
		user INT REFERENCES Users(id) ON DELETE CASCADE,
		role VARCHAR(255) NOT NULL UNIQUE,
		about TEXT NOT NULL,
		views INT,
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) GetProfile(username string) (*data.Profile, error) {

	user_id, f_err := s.getUserID(username)
	if f_err != nil {
		return nil, nil
	}

	profile := new(data.Profile)
	q := s.db.QueryRow("SELECT user, role, about, views FROM Profiles WHERE user = $1", user_id)
	err := q.Scan(
		&profile.User,
		&profile.Role,
		&profile.About,
		&profile.Views,
	)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *PostgresStorage) GetProfileByRole(role string) ([]*data.Profile, error) {
	rows, err := s.db.Query("SELECT user, role, views FROM Profiles WHERE role = $1", role)
	if err != nil {
		return nil, err
	}

	var profiles []*data.Profile
	for rows.Next() {
		profile := new(data.Profile)
		err := rows.Scan(
			profile.User,
			profile.Role,
			profile.About,
			profile.Views,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (s *PostgresStorage) CreateProfile(p data.Profile) error {

	user_id, f_err := s.getUserID(p.User.Username)
	if f_err != nil {
		return nil
	}

	query := `INSERT INTO Profiles (user, role, about, views)
	VALUES ($1, $2, $3);`

	_, err := s.db.Query(
		query,
		user_id,
		p.Role,
		p.About,
		p.Views,
	)
	return err
}

func (s *PostgresStorage) DeleteProfile(p data.Profile) error {
	q := "DELETE FROM Profiles WHERE User = $1 AND role = $2"

	user_id, f_err := s.getUserID(p.User.Username)
	if f_err != nil {
		return nil
	}

	_, err := s.db.Exec(q, user_id, p.Role)
	return err
}

// Profile

// Session

func (s *PostgresStorage) createSessionTable() error {
	query := `CREATE TABLE IF NOT EXISTS Sessions (
		id SERIAL PRIMARY KEY,
		user INT REFERENCES Users(id) ON DELETE CASCADE,
		key VARCHAR(255) NOT NULL UNIQUE,
		expires_on TIMESTAMP,
		expired BOOLEAN DEFAULT FALSE
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateSession(session data.Session) error {
	query := "INSERT INTO Sessions (user, key, expires) VALUES ($1, $2, $3)"

	user_id, f_err := s.getUserID(session.User.Username)
	if f_err != nil {
		return nil
	}

	_, err := s.db.Query(query, user_id, session.Key)
	return err
}

func (s *PostgresStorage) GetSession(u data.User) (*data.Session, error) {

	user_id, f_err := s.getUserID(u.Username)
	if f_err != nil {
		return nil, nil
	}

	session := new(data.Session)
	q := s.db.QueryRow("SELECT Key, expires_on FROM Sessions WHERE user = $1", user_id)
	err := q.Scan(
		&session.Key,
		&session.Expires_on,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *PostgresStorage) DeleteSession(ss data.Session) error {
	q := "DELETE FROM Session WHERE User = $1 AND expired = True"

	user_id, f_err := s.getUserID(ss.User.Username)
	if f_err != nil {
		return nil
	}

	_, err := s.db.Exec(q, user_id)
	return err
}

func (s *PostgresStorage) CancelSession(u data.User) error {

	user_id, f_err := s.getUserID(u.Username)
	if f_err != nil {
		return nil
	}

	q := "UPDATE Sessions SET expired = $1 WHERE id = $2"

	_, err := s.db.Exec(q, true, user_id)
	return err
}

// Session
