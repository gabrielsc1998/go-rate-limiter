package mysql

import "testing"

func TestNewMySQLDBConnection(t *testing.T) {
	db := NewMySQLDBConnection()

	if db == nil {
		t.Errorf("Expected db to be not nil")
	}
}

func TestConnect(t *testing.T) {
	db := NewMySQLDBConnection()
	options := MySQLConnectionOptions{
		Host:     "mysql",
		Port:     "3306",
		User:     "root",
		Password: "root",
		Database: "rate_limiter",
	}

	err := db.Connect(options)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClose(t *testing.T) {
	db := NewMySQLDBConnection()
	options := MySQLConnectionOptions{
		Host:     "mysql",
		Port:     "3306",
		User:     "root",
		Password: "root",
		Database: "rate_limiter",
	}

	err := db.Connect(options)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
