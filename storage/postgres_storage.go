package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/phillipmugisa/go_resume_generator/data"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	db, err := initPostgresDB()
	if err != nil {
		return nil, err
	}
	fmt.Println("Database Connection Successful...")
	return &PostgresStorage{
		db: db,
	}, nil
}

func initPostgresDB() (*sql.DB, error) {
	HOST := os.Getenv("POSTGRES_HOST")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DB")
	PORT := os.Getenv("POSTGRES_PORT")
	username := os.Getenv("POSTGRES_USER")

	// make db connection
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, username, password, database)

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, errors.New("couldnot connect to database")
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

	if err := s.createUserImageTable(); err != nil {
		return err
	}

	if err := s.createProfileTable(); err != nil {
		return err
	}

	if err := s.createSessionTable(); err != nil {
		return err
	}

	if err := s.createProjectTable(); err != nil {
		return err
	}

	if err := s.createEmploymentTable(); err != nil {
		return err
	}

	if err := s.createTechStackTable(); err != nil {
		return err
	}

	if err := s.createProjectTechStackTable(); err != nil {
		return err
	}

	if err := s.createEmploymentTechStackTable(); err != nil {
		return err
	}

	if err := s.createHobbiesTable(); err != nil {
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
		start_date TIMESTAMP NOT NULL,
		years_of_work INTEGER NULL,
		created_on TIMESTAMP NOT NULL,
		updated_on TIMESTAMP NOT NULL,
		last_sign_in TIMESTAMP NULL,
		portfolio VARCHAR(255) NULL,
		github VARCHAR(255) NULL,
		linkedin VARCHAR(255) NULL,
		twitter VARCHAR(255) NULL
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) createUserImageTable() error {
	query := `CREATE TABLE IF NOT EXISTS UserImages (
		id SERIAL PRIMARY KEY,
		filename VARCHAR(255) NOT NULL,
		user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) GetUsers(keys map[string]string) ([]*data.User, error) {

	query := "SELECT id, username, firstname, lastname, email, bio, phone, country, password FROM Users WHERE"
	if len(keys) > 0 {
		loop_counter := 0
		for k, v := range keys {

			var query_part string
			_, err := strconv.Atoi(v)
			if err != nil {
				query_part = fmt.Sprintf(" %s = '%s' ", k, v)
			} else {
				query_part = fmt.Sprintf(" %s = %s ", k, v)
			}

			if loop_counter > 0 {
				query = query + fmt.Sprintf(" AND %s ", query_part)
			} else {
				query = query + query_part
			}
			loop_counter++
		}
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	return scanUsers(rows)
}
func (s *PostgresStorage) CreateUser(u data.User) error {

	query := `INSERT INTO Users (username, firstname, lastname, email, password, phone, country, created_on, updated_on, bio, start_date, years_of_work)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

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
		u.Start_date,
		u.Years_of_work,
	)
	return err
}

func (s *PostgresStorage) CreateUserimage(u data.User, filaname string) error {

	user_id, f_err := s.getUserID(u.Username)
	if f_err != nil {
		return f_err
	}

	query := `INSERT INTO UserImages (filename, user_id) VALUES ($1, $2);`

	_, err := s.db.Query(
		query,
		filaname,
		user_id,
	)
	return err
}

func (s *PostgresStorage) SetUserSocials(u data.User) error {

	user_id, f_err := s.getUserID(u.Username)
	if f_err != nil {
		return f_err
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
		user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
		role VARCHAR(255) NOT NULL UNIQUE,
		about TEXT NOT NULL,
		views INTEGER
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) GetProfile(username string) (*data.Profile, error) {

	user_id, f_err := s.getUserID(username)
	if f_err != nil {
		return nil, f_err
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
		return f_err
	}

	query := `INSERT INTO Profiles (user_id, role, about, views)
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
		return f_err
	}

	_, err := s.db.Exec(q, user_id, p.Role)
	return err
}

// Profile

// Session

func (s *PostgresStorage) createSessionTable() error {
	query := `CREATE TABLE IF NOT EXISTS Sessions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
		key VARCHAR(255) NOT NULL UNIQUE,
		expires_on TIMESTAMP,
		expired BOOLEAN DEFAULT FALSE
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateSession(session data.Session) error {
	user_id, f_err := s.getUserID(session.User.Username)
	if f_err != nil {
		return f_err
	}

	// delete existing active sessions
	delete_query := "DELETE FROM Sessions WHERE user_id = $1"
	_, delete_err := s.db.Exec(delete_query, user_id)
	if delete_err != nil {
		return delete_err
	}

	// create new session
	query := "INSERT INTO Sessions (user_id, key, expires_on, expired) VALUES ($1, $2, $3, $4)"

	_, err := s.db.Query(query, user_id, session.Key, session.Expires_on, session.Expired)
	return err
}

func (s *PostgresStorage) GetSession(key string) (*data.Session, error) {
	var user_id int

	session := new(data.Session)
	q := s.db.QueryRow("SELECT * FROM Sessions WHERE Key = $1", key)
	err := q.Scan(
		&session.Id,
		&user_id,
		&session.Key,
		&session.Expires_on,
		&session.Expired,
	)

	if err != nil {
		return nil, err
	}

	users, err := s.GetUsers(map[string]string{"id": fmt.Sprintf("%d", user_id)})
	if err != nil {
		return nil, err
	}
	session.User = *users[0]
	return session, nil
}

func (s *PostgresStorage) DeleteSession(ss data.Session) error {
	q := "DELETE FROM Session WHERE User = $1 AND expired = True"

	user_id, f_err := s.getUserID(ss.User.Username)
	if f_err != nil {
		return f_err
	}

	_, err := s.db.Exec(q, user_id)
	return err
}

func (s *PostgresStorage) CancelSession(session data.Session) error {
	q := "UPDATE Sessions SET expired = $1 WHERE key = $2"

	_, err := s.db.Exec(q, true, session.Key)
	return err
}

// Session

// Projects
func (s *PostgresStorage) createProjectTable() error {
	query := `CREATE TABLE IF NOT EXISTS Projects (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		duration VARCHAR(255) NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NULL,
		status VARCHAR(255) DEFAULT 'COMPLETED',
		github VARCHAR(255) NULL,
		prod_link VARCHAR(255) NULL,
		description VARCHAR(255) NOT NULL,
		created_on TIMESTAMP,
		updated_on TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateProject(p data.Project) error {
	q := `INSERT INTO Projects (user_id, name, duration, start_date, end_date, status, github, prod_link, description, created_on, updated_on) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	user_id, f_err := s.getUserID(p.User.Username)
	if f_err != nil {
		return f_err
	}

	_, err := s.db.Query(q, user_id, p.Name, p.Duration, p.Start_date, p.End_date, p.Status, p.Github, p.Prod_link, p.Description, p.Created_on, p.Updated_on)
	return err
}

func (s *PostgresStorage) GetProjects(keys map[string]string) ([]*data.Project, error) {
	// expects a map with keys: id, username, name, stack
	// combines keys using and statement

	if len(keys) == 0 {
		return nil, errors.New("provide search keyword")
	}
	query := "SELECT * FROM Projects WHERE "

	loop_counter := 0
	for k, v := range keys {
		if k == "stack" {
			continue
		}
		if k == "username" {
			user_id, err := s.getUserID(v)
			if err != nil {
				continue
			}
			k = "user"
			v = fmt.Sprintf("%d", user_id)
		}

		var query_part string
		_, err := strconv.Atoi(v)
		if err != nil {
			query_part = fmt.Sprintf(" %s = '%s' ", k, v)
		} else {
			query_part = fmt.Sprintf(" %s = %s ", k, v)
		}

		if loop_counter > 0 {
			query = query + fmt.Sprintf(" AND %s ", query_part)
		} else {
			query = query + query_part
		}

		loop_counter++
	}

	projects := []*data.Project{}
	if loop_counter > 0 {
		rows, err := s.db.Query(query)
		if err != nil {
			return nil, err
		}

		projects, _ = s.scanProjects(rows)
	}

	stack_name, ok := keys["stack"]
	if ok {
		keywords := map[string]string{
			"name": stack_name,
		}
		p_id, ok := keys["id"]
		if ok {
			keywords["project_id"] = p_id
		}

		username, ok := keys["username"]
		if ok {
			keywords["username"] = username
		}

		p, err := s.GetProjectsByTechStack(keywords)
		if err != err {
			return nil, nil
		}

		// only add record that havent been added yet
	outer:
		for _, v := range p {
			for _, k := range projects {
				if k.Id == v.Id {
					// move to the next item in p (the outer loop)
					continue outer
				}
			}
			projects = append(projects, v)
		}
	}

	return projects, nil
}

func (s *PostgresStorage) GetProjectsByTechStack(keys map[string]string) ([]*data.Project, error) {
	// expects a map with keys: techstack_id, name, project_id, username

	if len(keys) == 0 {
		return nil, errors.New("provide search keyword")
	}

	query := "SELECT * FROM ProjectTechStacks WHERE "

	loop_counter := 0
	for k, v := range keys {
		if k == "name" {
			var stack_id int
			f_err := s.db.QueryRow("SELECT id FROM TechStacks WHERE name = $1", v).Scan(&stack_id)
			if f_err != nil {
				return nil, f_err
			}
			k = "techstack_id"
			v = fmt.Sprintf("%d", stack_id)
		}
		if k == "username" {
			continue
		}

		var query_part string
		_, err := strconv.Atoi(v)
		if err != nil {
			query_part = fmt.Sprintf(" %s = '%s' ", k, v)
		} else {
			query_part = fmt.Sprintf(" %s = %s ", k, v)
		}

		if loop_counter > 0 {
			query = query + fmt.Sprintf(" AND %s ", query_part)
		} else {
			query = query + query_part
		}

		loop_counter++
	}

	var projects []*data.Project

	if loop_counter == 0 {
		return nil, nil
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		record := new(struct {
			id           string
			techstack_id string
			project_id   string
		})
		err := rows.Scan(record)
		if err != nil {
			continue
		}

		// get projects with found id
		p, e := s.GetProjects(map[string]string{"id": record.project_id})
		if e != nil {
			return nil, e
		}
		username, ok := keys["username"]
		if ok {
			if p[0].User.Username != username {
				continue
			}
		}
		projects = append(projects, p[0])
	}

	return projects, nil
}

func (s *PostgresStorage) scanProjects(rows *sql.Rows) ([]*data.Project, error) {
	var projects []*data.Project
	for rows.Next() {
		project := new(data.Project)
		var user_id int
		err := rows.Scan(
			project.Id,
			user_id,
			project.Name,
			project.Duration,
			project.Start_date,
			project.End_date,
			project.Status,
			project.Github,
			project.Prod_link,
			project.Description,
			project.Created_on,
			project.Updated_on,
		)
		if err != nil {
			return projects, err
		}

		users, _ := s.GetUsers(map[string]string{"id": fmt.Sprintf("%d", user_id)})
		project.User = *users[0]

		projects = append(projects, project)
	}
	return projects, nil
}

func (s *PostgresStorage) DeleteProject(id int) error {
	q := "DELETE FROM Projects WHERE id = $1"

	_, err := s.db.Exec(q, id)
	return err
}

// Employment
func (s *PostgresStorage) createEmploymentTable() error {
	query := `CREATE TABLE IF NOT EXISTS Employments (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		employee VARCHAR(255) NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NULL,
		status VARCHAR(255) DEFAULT 'Current',
		prod_link VARCHAR(255) NULL,
		duration VARCHAR(255) NULL,
		description VARCHAR(255) NOT NULL,
		created_on TIMESTAMP,
		updated_on TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateEmployment(e data.Employment) error {
	q := `INSERT INTO Employment (user_id, name, employee, start_date, end_date, status, prod_link, duration, description, created_on, updated_on) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	user_id, f_err := s.getUserID(e.User.Username)
	if f_err != nil {
		return f_err
	}

	_, err := s.db.Query(q, user_id, e.Name, e.Employee, e.Start_date, e.End_date, e.Status, e.Prod_link, e.Duration, e.Description, e.Created_on, e.Updated_on)
	return err
}

func (s *PostgresStorage) GetEmployments(keys map[string]string) ([]*data.Employment, error) {
	// expects a map with keys: id, username, name, stack
	// combines keys using and statement

	if len(keys) == 0 {
		return nil, errors.New("provide search keyword")
	}
	query := "SELECT * FROM Employments WHERE "

	loop_counter := 0
	for k, v := range keys {
		if k == "stack" {
			continue
		}
		if k == "username" {
			user_id, err := s.getUserID(v)
			if err != nil {
				continue
			}
			k = "user"
			v = fmt.Sprintf("%d", user_id)
		}

		var query_part string
		_, err := strconv.Atoi(v)
		if err != nil {
			query_part = fmt.Sprintf(" %s = '%s' ", k, v)
		} else {
			query_part = fmt.Sprintf(" %s = %s ", k, v)
		}

		if loop_counter > 0 {
			query = query + fmt.Sprintf(" AND %s ", query_part)
		} else {
			query = query + query_part
		}

		loop_counter++
	}

	employments := []*data.Employment{}
	if loop_counter > 0 {
		rows, err := s.db.Query(query)
		if err != nil {
			return nil, err
		}

		employments, _ = s.scanEmployments(rows)
	}

	stack_name, ok := keys["stack"]
	if ok {
		keywords := map[string]string{
			"name": stack_name,
		}
		p_id, ok := keys["id"]
		if ok {
			keywords["project_id"] = p_id
		}

		p, err := s.GetEmploymentsByTechStack(keywords)
		if err != err {
			return nil, nil
		}

		// only add record that havent been added yet
	outer:
		for _, v := range p {
			for _, k := range employments {
				if k.Id == v.Id {
					// move to the next item in p (the outer loop)
					continue outer
				}
			}
			employments = append(employments, v)
		}
	}

	return employments, nil
}

func (s *PostgresStorage) GetEmploymentsByTechStack(keys map[string]string) ([]*data.Employment, error) {
	// expects a map with keys: id, name, employment_id, username

	if len(keys) == 0 {
		return nil, errors.New("provide search keyword")
	}

	query := "SELECT * FROM EmploymentTechStacks WHERE "

	loop_counter := 0
	for k, v := range keys {
		if k == "name" {
			var stack_id int
			f_err := s.db.QueryRow("SELECT id FROM TechStacks WHERE name = $1", v).Scan(&stack_id)
			if f_err != nil {
				return nil, f_err
			}
			k = "techstack_id"
			v = fmt.Sprintf("%d", stack_id)
		}
		if k == "username" {
			continue
		}

		var query_part string
		_, err := strconv.Atoi(v)
		if err != nil {
			query_part = fmt.Sprintf(" %s = '%s' ", k, v)
		} else {
			query_part = fmt.Sprintf(" %s = %s ", k, v)
		}

		if loop_counter > 0 {
			query = query + fmt.Sprintf(" AND %s ", query_part)
		} else {
			query = query + query_part
		}

		loop_counter++
	}

	if loop_counter == 0 {
		return nil, nil
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	var employments []*data.Employment
	for rows.Next() {
		record := new(struct {
			id            string
			techstack_id  string
			employment_id string
		})
		err := rows.Scan(record)
		if err != nil {
			continue
		}

		// get employments with found id
		p, e := s.GetEmployments(map[string]string{"id": record.employment_id})
		if e != nil {
			return nil, e
		}

		username, ok := keys["username"]
		if ok {
			// if username was passed, only add employments with this stack and user
			if p[0].User.Username != username {
				continue
			}
		}
		employments = append(employments, p[0])
	}

	return employments, nil
}

func (s *PostgresStorage) scanEmployments(rows *sql.Rows) ([]*data.Employment, error) {
	var Employments []*data.Employment
	for rows.Next() {
		Employment := new(data.Employment)
		var user_id int
		err := rows.Scan(
			Employment.Id,
			user_id,
			Employment.Name,
			Employment.Employee,
			Employment.Start_date,
			Employment.End_date,
			Employment.Status,
			Employment.Prod_link,
			Employment.Duration,
			Employment.Description,
			Employment.Created_on,
			Employment.Updated_on,
		)
		if err != nil {
			return Employments, err
		}

		users, _ := s.GetUsers(map[string]string{"id": fmt.Sprintf("%d", user_id)})
		Employment.User = *users[0]

		Employments = append(Employments, Employment)
	}
	return Employments, nil
}

func (s *PostgresStorage) DeleteEmployment(id int) error {
	q := "DELETE FROM Employments WHERE id = $1"

	_, err := s.db.Exec(q, id)
	return err
}

// // Hobby
func (s *PostgresStorage) createHobbiesTable() error {
	query := `CREATE TABLE IF NOT EXISTS Hobbies (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateHobby(h data.Hobby) error {
	q := `INSERT INTO Hobbies (user_id, name) VALUES ($1, $2)`

	user_id, f_err := s.getUserID(h.User.Username)
	if f_err != nil {
		return f_err
	}

	_, err := s.db.Query(q, user_id, h.Name)
	return err

}

func (s *PostgresStorage) GetHobbies(keys map[string]string) ([]*data.Hobby, error) {
	// expects keys: id, username, name

	if len(keys) == 0 {
		return nil, errors.New("provide search keyword")
	}

	query := "SELECT * FROM Hobbies WHERE "

	loop_counter := 0
	for k, v := range keys {
		if k == "username" {
			user_id, err := s.getUserID(v)
			if err != nil {
				continue
			}
			k = "user"
			v = fmt.Sprintf("%d", user_id)
		}

		var query_part string
		_, err := strconv.Atoi(v)
		if err != nil {
			query_part = fmt.Sprintf(" %s = '%s' ", k, v)
		} else {
			query_part = fmt.Sprintf(" %s = %s ", k, v)
		}

		if loop_counter > 0 {
			query = query + fmt.Sprintf(" AND %s ", query_part)
		} else {
			query = query + query_part
		}

		loop_counter++
	}

	if loop_counter == 0 {
		return nil, nil
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	var hobbies []*data.Hobby
	for rows.Next() {
		hobbie := new(data.Hobby)
		err := rows.Scan(hobbie)
		if err != nil {
			continue
		}
		hobbies = append(hobbies, hobbie)
	}
	return hobbies, nil
}

func (s *PostgresStorage) DeleteHobby(id int) error {
	q := "DELETE FROM Hobbies WHERE id = $1"

	_, err := s.db.Exec(q, id)
	return err
}

// // TechStack
func (s *PostgresStorage) createTechStackTable() error {
	query := `CREATE TABLE IF NOT EXISTS TechStacks (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL
	)`
	_, err := s.db.Exec(query)
	return err
}

// TECH STACK RELATIONSHIPS (M:M)
func (s *PostgresStorage) createProjectTechStackTable() error {
	query := `CREATE TABLE IF NOT EXISTS ProjectTechStacks (
		id SERIAL PRIMARY KEY,
		techstack_id INTEGER REFERENCES TechStacks(id) ON DELETE CASCADE,
		project_id INTEGER REFERENCES Projects(id) ON DELETE CASCADE
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) createEmploymentTechStackTable() error {
	query := `CREATE TABLE IF NOT EXISTS EmploymentTechStacks (
		id SERIAL PRIMARY KEY,
		techstack_id INTEGER REFERENCES TechStacks(id) ON DELETE CASCADE,
		employment_id INTEGER REFERENCES Employments(id) ON DELETE CASCADE
	)`
	_, err := s.db.Exec(query)
	return err
}

// TECH STACK RELATIONSHIPS

func (s *PostgresStorage) CreateTechStack(t data.TechStack) error {
	q := `INSERT INTO TechStacks (user_id, name) VALUES ($1, $2)`

	user_id, f_err := s.getUserID(t.User.Username)
	if f_err != nil {
		return f_err
	}

	_, err := s.db.Query(q, user_id, t.Name)
	return err
}

func (s *PostgresStorage) GetTechStacks(keys map[string]string) ([]*data.TechStack, error) {
	// expects keys: id, username, name

	if len(keys) == 0 {
		return nil, errors.New("provide search keyword")
	}

	query := "SELECT * FROM TechStacks WHERE "

	loop_counter := 0
	for k, v := range keys {
		if k == "username" {
			user_id, err := s.getUserID(v)
			if err != nil {
				continue
			}
			k = "user"
			v = fmt.Sprintf("%d", user_id)
		}

		var query_part string
		_, err := strconv.Atoi(v)
		if err != nil {
			query_part = fmt.Sprintf(" %s = '%s' ", k, v)
		} else {
			query_part = fmt.Sprintf(" %s = %s ", k, v)
		}

		if loop_counter > 0 {
			query = query + fmt.Sprintf(" AND %s ", query_part)
		} else {
			query = query + query_part
		}

		loop_counter++
	}

	if loop_counter == 0 {
		return nil, nil
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	return s.scanTechStack(rows)
}

func (s *PostgresStorage) scanTechStack(rows *sql.Rows) ([]*data.TechStack, error) {
	var stacks []*data.TechStack
	for rows.Next() {
		stack := new(data.TechStack)
		err := rows.Scan(stack)
		if err != nil {
			continue
		}
		stacks = append(stacks, stack)
	}
	return stacks, nil
}

// returns techstacks for a given project
func (s *PostgresStorage) GetProjectTechStacks(keys map[string]string) ([]*data.TechStack, error) {
	// expects keys: project_id, project_name

	if len(keys) == 0 {
		return nil, errors.New("provide search keyword")
	}

	query := "SELECT * FROM ProjectTechStacks WHERE "

	project_id, _ := keys["project_id"]

	project_name, ok := keys["project_name"]
	if ok {
		// get project id
		results, err := s.GetProjects(map[string]string{"name": project_name})
		if err != nil {
			return nil, err
		}
		project_id = fmt.Sprintf("%d", results[0].Id)
	}

	rows, err := s.db.Query(fmt.Sprintf("%s project_id = %s", query, project_id))
	if err != nil {
		return nil, err
	}

	var stacks []*data.TechStack

	for rows.Next() {
		record := new(struct {
			id           string
			techstack_id string
			project_id   string
		})
		err := rows.Scan(record)
		if err != nil {
			continue
		}

		// get techstacks with found id
		stack, err := s.GetTechStacks(map[string]string{"id": record.techstack_id})
		if err != nil {
			return nil, err
		}
		username, ok := keys["username"]
		if ok {
			if stack[0].User.Username != username {
				continue
			}
		}
		stacks = append(stacks, stack[0])
	}

	return stacks, nil

}

// returns techstacks for a given employment
func (s *PostgresStorage) GetEmploymentTechStacks(keys map[string]string) ([]*data.TechStack, error) {
	// expects keys: employment_id, employment_name

	if len(keys) == 0 {
		return nil, errors.New("provide search keyword")
	}

	query := "SELECT * FROM EmploymentTechStacks WHERE "

	employment_id, _ := keys["employment_id"]

	// get id for given username
	employment_name, ok := keys["employment_name"]
	if ok {
		// get employment id
		results, err := s.GetEmployments(map[string]string{"name": employment_name})
		if err != nil {
			return nil, err
		}
		employment_id = fmt.Sprintf("%d", results[0].Id)
	}

	rows, err := s.db.Query(fmt.Sprintf("%s employment_id = %s", query, employment_id))
	if err != nil {
		return nil, err
	}

	var stacks []*data.TechStack

	for rows.Next() {
		record := new(struct {
			id            string
			techstack_id  string
			employment_id string
		})
		err := rows.Scan(record)
		if err != nil {
			continue
		}

		// get techstacks with found id
		stack, e := s.GetTechStacks(map[string]string{"id": record.techstack_id})
		if e != nil {
			return nil, e
		}
		username, ok := keys["username"]
		if ok {
			if stack[0].User.Username != username {
				continue
			}
		}
		stacks = append(stacks, stack[0])
	}

	return stacks, nil

}

func (s *PostgresStorage) DeleteTechStack(id int) error {
	q := "DELETE FROM TechStacks WHERE id = $1"

	_, err := s.db.Exec(q, id)
	return err
}

func (s *PostgresStorage) AddTechStackToProject(t data.TechStack, p data.Project) error {
	// fetch project
	projects, err := s.GetProjects(map[string]string{"name": p.Name, "id": fmt.Sprint(p.Id)})
	if err != nil {
		return err
	}

	// fetch tech stack
	stacks, err := s.GetTechStacks(map[string]string{"name": p.Name, "id": fmt.Sprint(p.Id)})
	if err != nil {
		return err
	}

	// save to db
	_, write_err := s.db.Exec("INSERT INTO ProjectTechStacks (techstack_id, project_id) VALUES ($1, $2)", stacks[0].Id, projects[0].Id)
	return write_err
}

func (s *PostgresStorage) AddTechStackToEmployment(t data.TechStack, p data.Project) error {
	// fetch employment
	employments, err := s.GetEmployments(map[string]string{"name": p.Name, "id": fmt.Sprint(p.Id)})
	if err != nil {
		return err
	}

	// fetch tech stack
	stacks, err := s.GetTechStacks(map[string]string{"name": p.Name, "id": fmt.Sprint(p.Id)})
	if err != nil {
		return err
	}

	// save to db
	_, write_err := s.db.Exec("INSERT INTO ProjectTechStacks (techstack_id, project_id) VALUES ($1, $2)", stacks[0].Id, employments[0].Id)
	return write_err
}
