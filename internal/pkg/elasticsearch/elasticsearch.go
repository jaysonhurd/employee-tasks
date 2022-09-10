package elasticsearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/jaysonhurd/employee-tasks/internal/pkg/postgres"
	"github.com/jaysonhurd/employee-tasks/pkg/tasks/models"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type ElasticSearcher interface {
	LoadEmployeesFromPostgres() error
	LoadTasksFromPostgres() error
	EmployeeByNickname(nickname string) (employee []byte, err error)
	EmptyES() error
	AllEmployees() (employee []models.Employee, err error)
	AllTasks() (tasks []models.Task, err error)
}

const (
	esCertPath = "config/http_ca.crt"
)

type elasticSearch struct {
	esconn *elasticsearch.Client
	pgconn postgres.Postgreser
	l      *zerolog.Logger
}

func New(
	esconn *elasticsearch.Client,
	pgconn postgres.Postgreser,
	l *zerolog.Logger,
) ElasticSearcher {
	return &elasticSearch{
		esconn: esconn,
		pgconn: pgconn,
		l:      l,
	}
}

func (e *elasticSearch) AllEmployees() (employees []models.Employee, err error) {
	var buffer bytes.Buffer
	var result map[string]interface{}

	response, err := e.esconn.Search(e.esconn.Search.WithIndex("employees"), e.esconn.Search.WithBody(&buffer))
	if err != nil {
		return employees, err
	}

	json.NewDecoder(response.Body).Decode(&result)

	count := len(result["hits"].(map[string]interface{})["hits"].([]interface{}))

	for i := 0; i < count; i++ {
		e1 := result["hits"].(map[string]interface{})["hits"].([]interface{})[i]
		e2 := e1.(map[string]interface{})
		employee := e2["_source"].(map[string]interface{})

		employeeID, err := strconv.Atoi(fmt.Sprintf("%v", employee["employee_id"]))
		if err != nil {
			return employees, err
		}
		nickname := employee["nickname"]
		firstName := employee["first_name"]
		lastName := employee["last_name"]
		streetAddress := employee["street_address"]
		city := employee["city"]
		state := employee["state"]
		zip := employee["state"]

		employeeFinal := models.Employee{
			Employee_id:    employeeID,
			Nickname:       fmt.Sprintf("%v", nickname),
			First_name:     fmt.Sprintf("%v", firstName),
			Last_name:      fmt.Sprintf("%v", lastName),
			Street_address: fmt.Sprintf("%v", streetAddress),
			City:           fmt.Sprintf("%v", city),
			State:          fmt.Sprintf("%v", state),
			Zip:            fmt.Sprintf("%v", zip),
		}
		employees = append(employees, employeeFinal)
	}

	return employees, err
}

func (e *elasticSearch) AllTasks() (tasks []models.Task, err error) {
	var buffer bytes.Buffer
	var result map[string]interface{}

	response, err := e.esconn.Search(e.esconn.Search.WithIndex("tasks"), e.esconn.Search.WithBody(&buffer))
	if err != nil {
		return tasks, err
	}

	json.NewDecoder(response.Body).Decode(&result)

	// TODO: I know this is messy as all hell.  Map/string/interface.  Maybe I will get to know ElasticSearch better and make this more efficient
	count := len(result["hits"].(map[string]interface{})["hits"].([]interface{}))

	for i := 0; i < count; i++ {
		t1 := result["hits"].(map[string]interface{})["hits"].([]interface{})[i]
		t2 := t1.(map[string]interface{})
		task := t2["_source"].(map[string]interface{})

		taskID, err := strconv.Atoi(fmt.Sprintf("%v", task["id"]))
		if err != nil {
			return tasks, err
		}
		name := task["name"]
		description := task["description"]
		createTime := task["create_time"]
		// TODO: figure out how to convert this  []interface{} into an array.
		//owners := task["owners"]
		//value := reflect.ValueOf(owners)
		//fmt.Println(value)

		private, err := strconv.ParseBool(fmt.Sprintf("%v", task["private"]))
		if err != nil {
			return nil, err
		}
		dueBy := task["due_by"]

		taskFinal := models.Task{
			ID:          taskID,
			Name:        fmt.Sprintf("%v", name),
			Description: fmt.Sprintf("%v", description),
			Create_time: fmt.Sprintf("%v", createTime),
			Owners:      nil,
			Private:     private,
			Due_by:      fmt.Sprintf("%v", dueBy),
		}
		tasks = append(tasks, taskFinal)
	}

	return tasks, err
}

func (e *elasticSearch) EmployeeByNickname(nickname string) (employee []byte, err error) {
	var buffer bytes.Buffer
	var output []byte

	response, err := e.esconn.Search(e.esconn.Search.WithIndex("employee"), e.esconn.Search.WithBody(&buffer))
	if err != nil {
		return employee, err
	}
	fmt.Println(response.Body.Read(output))

	return employee, err
}

func (e *elasticSearch) LoadEmployeesFromPostgres() error {
	var employees []models.Employee
	ctx := context.Background()
	employees, err := e.pgconn.AllEmployees()
	if err != nil {
		return err
	}

	i := 0
	for range employees {
		employeeJSON, err := json.Marshal(employees[i])
		if err != nil {
			return err
		}
		req := esapi.IndexRequest{
			Index:      "employees",
			DocumentID: strconv.Itoa(i + 1),
			Body:       bytes.NewReader(employeeJSON),
			Refresh:    "true",
		}
		_, err = req.Do(ctx, e.esconn)
		if err != nil {
			return err
		}
		i++
	}

	return nil
}

func (e *elasticSearch) LoadTasksFromPostgres() error {
	var tasks []models.Task
	ctx := context.Background()
	tasks, err := e.pgconn.AllTasks()
	if err != nil {
		return err
	}

	i := 0
	for range tasks {
		taskJSON, err := json.Marshal(tasks[i])
		if err != nil {
			return err
		}
		req := esapi.IndexRequest{
			Index:      "tasks",
			DocumentID: strconv.Itoa(i + 1),
			Body:       bytes.NewReader(taskJSON),
			Refresh:    "true",
		}
		_, err = req.Do(ctx, e.esconn)
		if err != nil {
			return err
		}
		i++
	}

	return nil
}

func (e *elasticSearch) EmptyES() error {

	indices := make([]string, 2)
	indices[0] = "employee"
	indices[1] = "tasks"

	del := esapi.IndicesDeleteRequest{Index: indices}

	_, err := del.Do(context.Background(), e.esconn)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func NewESConnection(config elasticsearch.Config, l *zerolog.Logger) (*elasticsearch.Client, error) {
	var es *elasticsearch.Client
	var err error

	cert, err := getESCACert(esCertPath)
	if err != nil {
		return es, err
	}
	config.CACert = cert
	config.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	es, err = elasticsearch.NewClient(config)
	if err != nil {
		return es, err
	}

	es.Ping.WithHuman()
	l.Info().Msgf("elasticsearch connection established with " + config.Addresses[0])

	return es, nil
}

func getESCACert(path string) (cert []byte, err error) {

	_, err = os.Stat(path)
	if err != nil {
		return nil, err
	}
	cert, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
