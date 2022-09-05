package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strconv"
)

type Postgreser interface {
	AllTasks() ([]Task, error)
	EmployeeByID(ID string) (Employee, error)
	EmployeeByNickname(nickname string) (Employee, error)
	AllEmployees() ([]Employee, error)
}

type PostgresConn struct {
	ConnectString string
	SqlClient     *sql.DB
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"sslmode"`
}

type postgres struct {
	conn *PostgresConn
	l    *zerolog.Logger
}

type Task struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Create_time string        `json:"create_time"`
	Owners      pq.Int32Array `json:"owners"`
	Private     bool          `json:"private"`
	Due_by      string        `json:"due_by"`
}

type Employee struct {
	Employee_id    int    `json:"employee_id"`
	Nickname       string `json:"nickname"`
	First_name     string `json:"first_name"`
	Last_name      string `json:"last_name"`
	Street_address string `json:"street_address"`
	City           string `json:"city"`
	State          string `json:"state"`
	Zip            string `json:"zip"`
}

func New(
	conn *PostgresConn,
	l *zerolog.Logger,
) Postgreser {
	return &postgres{
		conn: conn,
		l:    l,
	}
}

func (p *postgres) AllTasks() ([]Task, error) {
	tasks := make([]Task, 0)
	db := p.conn.SqlClient
	query := `SELECT * FROM workers.public.tasks where private != true`

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Msgf(err.Error())
		return tasks, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return tasks, errors.New("Unable to run query to get all tasks")
	}
	defer rows.Close()

	for rows.Next() {
		var row Task
		if err = rows.Scan(&row.ID, &row.Name, &row.Description, &row.Create_time, &row.Owners, &row.Private, &row.Due_by); err != nil {
			p.l.Error().Msgf(fmt.Sprintf("Error scanning rows in All Task request to Postgres DB - %s", err.Error()))
		}
		tasks = append(tasks, Task{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
			Create_time: row.Create_time,
			Owners:      row.Owners,
			Private:     row.Private,
			Due_by:      row.Due_by,
		})
	}

	return tasks, nil
}

func (p *postgres) AllEmployees() ([]Employee, error) {
	employees := make([]Employee, 0)
	db := p.conn.SqlClient
	query := `SELECT * FROM workers.public.employees`

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error().Msgf(err.Error())
		return employees, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return employees, errors.New("Unable to run query to get all employees")
	}
	defer rows.Close()

	for rows.Next() {
		var row Employee
		if err = rows.Scan(&row.Employee_id, &row.Nickname, &row.First_name, &row.Last_name, &row.Street_address, &row.City, &row.State, &row.Zip); err != nil {
			p.l.Error().Msgf(fmt.Sprintf("Error scanning rows in All Task request to Postgres DB - %s", err.Error()))
		}
		employees = append(employees, Employee{
			Employee_id:    row.Employee_id,
			Nickname:       row.Nickname,
			First_name:     row.First_name,
			Last_name:      row.Last_name,
			Street_address: row.Street_address,
			City:           row.City,
			State:          row.State,
			Zip:            row.Zip,
		})
	}

	return employees, nil
}

func (p *postgres) EmployeeByID(ID string) (Employee, error) {
	var employee Employee
	var err error
	id, err := strconv.Atoi(ID)
	if err != nil {
		return employee, err
	}

	db := p.conn.SqlClient
	query := fmt.Sprintf(`SELECT * FROM workers.public.employees where employee_id = %d`, id)

	db.Prepare(query)

	row := db.QueryRow(query)
	err = row.Scan(&employee.Employee_id, &employee.Nickname, &employee.First_name, &employee.Last_name, &employee.Street_address, &employee.City, &employee.State, &employee.Zip)
	if err != nil {
		return employee, err
	}

	return employee, nil
}

func (p *postgres) EmployeeByNickname(nickname string) (Employee, error) {
	var employee Employee
	var err error
	if err != nil {
		return employee, err
	}

	db := p.conn.SqlClient
	query := fmt.Sprintf(`SELECT * FROM workers.public.employees where nickname = '%s'`, nickname)

	db.Prepare(query)

	row := db.QueryRow(query)
	err = row.Scan(&employee.Employee_id, &employee.Nickname, &employee.First_name, &employee.Last_name, &employee.Street_address, &employee.City, &employee.State, &employee.Zip)
	if err != nil {
		return employee, err
	}

	return employee, nil
}

// NewPostgresConnection - returns a PostgresConn connection and its conenct string for reference.  This can then be passed down
// through the service so that multiple connections do need to be made unncessarily.
func NewPostgresConnection(c PostgresConfig, l *zerolog.Logger) (db *PostgresConn, err error) {
	cString := fmt.Sprintf("host=%s port=%d database=%s user=%s password=%s sslmode=%s", c.Host, c.Port, c.Database, c.User, c.Password, c.SSLMode)
	m := &PostgresConn{
		ConnectString: cString,
	}
	m.SqlClient, err = sql.Open("postgres", cString)
	if err != nil {
		return nil, err
	}
	err = m.SqlClient.Ping()
	if err != nil {
		l.Error().Msgf("unable to connect to postgres")
		return
	}
	l.Error().Msgf("connection to PostgresConn works!")

	return m, nil
}
