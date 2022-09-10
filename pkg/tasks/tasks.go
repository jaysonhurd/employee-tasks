package tasks

import (
	"github.com/gin-gonic/gin"
	"github.com/jaysonhurd/employee-tasks/internal/pkg/elasticsearch"
	"github.com/jaysonhurd/employee-tasks/internal/pkg/postgres"
	"github.com/jaysonhurd/employee-tasks/pkg/tasks/models"
	"github.com/rs/zerolog"
	"net/http"
)

type Tasks interface {
	MakeRoutes() error
}

type tasks struct {
	postgres postgres.Postgreser
	elastic  elasticsearch.ElasticSearcher
	config   models.Config
	l        *zerolog.Logger
}

func New(
	p postgres.Postgreser,
	e elasticsearch.ElasticSearcher,
	c models.Config,
	l *zerolog.Logger,
) Tasks {
	return &tasks{
		postgres: p,
		elastic:  e,
		config:   c,
		l:        l,
	}
}

func (t *tasks) MakeRoutes() error {
	// Run gin endpoints
	r := gin.Default()
	rp := r.Group("/api/v1/")

	rp.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	rp.GET("/postgres/tasks", func(c *gin.Context) {
		tasks, err := t.postgres.AllTasks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, tasks)
	})

	rp.GET("/postgres/employee/id/:id", func(c *gin.Context) {
		id := c.Param("id")
		employee, err := t.postgres.EmployeeByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, &employee)
	})

	rp.GET("/postgres/employee/nickname/:nickname", func(c *gin.Context) {
		nickname := c.Param("nickname")
		employee, err := t.postgres.EmployeeByNickname(nickname)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, &employee)
	})

	rp.POST("/elasticsearch/load/employees", func(c *gin.Context) {
		err := t.elastic.LoadEmployeesFromPostgres()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "data loaded into elasticsearch successfully!"})
	})

	rp.POST("/elasticsearch/load/tasks", func(c *gin.Context) {
		err := t.elastic.LoadTasksFromPostgres()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "data loaded into elasticsearch successfully!"})
	})

	rp.GET("/elasticsearch/employees", func(c *gin.Context) {
		employees, err := t.elastic.AllEmployees()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, &employees)
	})

	rp.GET("/elasticsearch/tasks", func(c *gin.Context) {
		employees, err := t.elastic.AllTasks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, &employees)
	})

	err := r.Run("127.0.0.1:8080")

	return err
}
