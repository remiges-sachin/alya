package alya

import (
	"database/sql"
	"fmt"

	"github.com/remiges-tech/alya/config"
	"github.com/remiges-tech/alya/container"

	_ "github.com/lib/pq" // Assuming PostgreSQL, replace with your driver
)

// Queries represents the client for the database.
type Queries struct {
	db *sql.DB
}

// New creates a new Queries instance.
func New(db *sql.DB) *Queries {
	return &Queries{
		db: db,
	}
}

// SQLCService wraps the sqlc generated code and manages the database connection.
type SQLCService struct {
	DB        *sql.DB
	Queries   *Queries
	Container *container.Container
}

// NewSQLCService creates a new SQLCService with the necessary dependencies.
func NewSQLCService(container *container.Container) *SQLCService {
	return &SQLCService{
		DB:        nil, // DB will be initialized in Init
		Container: container,
	}
}

// Init establishes a new database connection and prepares the query struct.
func (s *SQLCService) Init() error {
	var err error

	// Retrieve the config object from the container
	configObj, err := s.Container.Resolve("config")
	if err != nil {
		return err
	}

	// Assert that the config object implements the config.Config interface
	config, ok := configObj.(config.Config)
	if !ok {
		return fmt.Errorf("config object does not implement the config.Config interface")
	}

	// Use the Get method to get the PostgreSQL connection parameters
	dsn, err := config.Get("postgres_dsn")
	if err != nil {
		return err
	}

	s.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	// Assuming you want to ping the database to ensure the connection is established
	if err = s.DB.Ping(); err != nil {
		return err
	}

	s.Queries = New(s.DB)
	return nil
}

// Close terminates the database connection.
func (s *SQLCService) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
