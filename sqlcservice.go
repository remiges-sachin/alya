package alya

import (
	"database/sql"

	"github.com/remiges-tech/alya/db"

	_ "github.com/lib/pq" // Assuming PostgreSQL, replace with your driver
)

// SQLCService wraps the sqlc generated code and manages the database connection.
type SQLCService struct {
	DB      *sql.DB
	Queries *db.Queries
}

// NewSQLCService creates a new SQLCService with the necessary dependencies.
func NewSQLCService(dataSourceName string) *SQLCService {
	return &SQLCService{
		DB: nil, // DB will be initialized in Init
	}
}

// Init establishes a new database connection and prepares the query struct.
func (s *SQLCService) Init() error {
	var err error
	s.DB, err = sql.Open("postgres", "your-database-connection-string") // Replace with your DSN
	if err != nil {
		return err
	}

	// Assuming you want to ping the database to ensure the connection is established
	if err = s.DB.Ping(); err != nil {
		return err
	}

	s.Queries = db.New(s.DB)
	return nil
}

// Close terminates the database connection.
func (s *SQLCService) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
