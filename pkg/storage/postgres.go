package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sergeikus/go-rest-template/pkg/types"
)

// PostgresStorage represents a Postgres database
type PostgresStorage struct {
	DSN     string
	PGXPool *pgxpool.Pool
}

// DefinePostgresStorage PostgresStorage fields
func DefinePostgresStorage(user, password, dbname, host string, port int) *PostgresStorage {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s", user, password, host, port, dbname)
	ps := &PostgresStorage{
		DSN: dsn,
	}
	return ps
}

// Connect performs connection to the database
func (ps *PostgresStorage) Connect() error {
	//   # Example DSN
	// user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10
	dbPool, err := pgxpool.Connect(context.Background(), ps.DSN)
	if err != nil {
		return fmt.Errorf("failed to perform database connection to '%s': %v", ps.DSN, err)
	}
	ps.PGXPool = dbPool
	return nil
}

// Close closes connection to the database
func (ps *PostgresStorage) Close() {
	ps.PGXPool.Close()
}

// Store performs storage of data in database, returns stored data auto generated primary key
func (ps *PostgresStorage) Store(data string) (id int, err error) {
	sql := `
	INSERT INTO data_table (string)
	VALUES ($1)
	RETURNING id
	`
	if err := ps.PGXPool.QueryRow(context.Background(), sql, data).Scan(&id); err != nil {
		return id, fmt.Errorf("failed to store data: %v", err)
	}
	return id, nil
}

// GetAll returns all data (rows) from 'data_table'
func (ps *PostgresStorage) GetAll() ([]types.Data, error) {
	sql := `
	SELECT * FROM data_table
	`
	rows, err := ps.PGXPool.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("failed to get all data from table: %v", err)
	}
	defer rows.Close()

	var result []types.Data
	for rows.Next() {
		var d types.Data
		rows.Scan(&d.ID, &d.String)
		result = append(result, d)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("encountered an error while reading rows: %v", err)
	}

	return result, nil
}

// GetKey returns data for a particular key
func (ps *PostgresStorage) GetKey(key int) (types.Data, error) {
	sql := `
	SELECT * FROM data_table 
	WHERE id=$1
	`
	var d types.Data
	if err := ps.PGXPool.QueryRow(context.Background(), sql, key).Scan(&d.ID, &d.String); err != nil {
		return d, fmt.Errorf("failed to query data for '%d' key: %v", key, err)
	}
	return d, nil
}
