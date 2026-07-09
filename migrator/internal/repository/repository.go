// TODO: Настроить context
// TODO: Добавить интерфейсы

package repository

import (
	"context"
	"fmt"
	"log"
	ewrap "postgres-migrator/internal/lib"

	"github.com/jackc/pgx/v5"
)

const (
	CREATE_TABLE_QUERY = `CREATE TABLE IF NOT EXISTS migrations (
							id SERIAL PRIMARY KEY,
							version BIGINT UNIQUE NOT NULL,
							name VARCHAR(100),
							applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	INSERT_MIGRATION_QUERY = `INSERT INTO migrations (version, name) VALUES ($1, $2);`

	MIGRATIONS_QUERY = `SELECT version FROM migrations;`
)

type Credentials struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB_name  string `json:"postgres_db_name"`
}

type Connection struct {
	conn *pgx.Conn
}

func NewConnection(creds Credentials) (*Connection, error) {
	ctx := context.Background()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", creds.Username, creds.Password, creds.Host, creds.Port, creds.DB_name)

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't connect to database", err)
	}

	c := Connection{
		conn: conn,
	}

	if err = c.setupTable(); err != nil {
		return nil, ewrap.Wrap("Couldn't setup table", err)
	}

	return &c, nil
}

func (c *Connection) setupTable() error {
	ctx := context.Background()

	_, err := c.conn.Exec(ctx, CREATE_TABLE_QUERY)

	if err != nil {
		return ewrap.Wrap("Couldn't check existence of table", err)
	}

	return nil
}

func (c *Connection) Close() error {
	ctx := context.Background()
	return c.conn.Close(ctx)
}

func (c *Connection) ApplyMigration(version uint64, name string, content string) error {
	ctx := context.Background()

	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return ewrap.Wrap("Couldn't begin transaction", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, content)
	if err != nil {
		return ewrap.Wrap("Couldn't execute content", err)
	}

	_, err = tx.Exec(ctx, INSERT_MIGRATION_QUERY, version, name)
	if err != nil {
		return ewrap.Wrap("Couldn't insert migration", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return ewrap.Wrap("Couldn't finish transaction", err)
	}

	return nil
}

func (c *Connection) Migrations() (map[uint64]bool, error) {
	ctx := context.Background()

	rows, err := c.conn.Query(ctx, MIGRATIONS_QUERY)
	if err != nil {
		return nil, ewrap.Wrap("Couldn't get migrations", err)
	}

	result := make(map[uint64]bool, 5)
	for rows.Next() {
		var val uint64
		if err := rows.Scan(&val); err != nil {
			return nil, ewrap.Wrap("Couldn't scan value", err)
		}
		result[val] = true
		log.Println("Existing migration: ", val)
	}

	if err := rows.Err(); err != nil {
		return nil, ewrap.Wrap("Couldn't iterate migrations", err)
	}

	return result, nil
}
