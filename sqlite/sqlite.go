package sqlite

import (
	"context"
	"strings"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

// Reason behind this is that
// I don't want to import two packages that
// starts with sqlite into my project
// I can use the type alias in my project
type Conn = sqlite.Conn
type Stmt = sqlite.Stmt

var SavePoint = sqlitex.Save

// Database struct which holds pool of connection
type Database struct {
	pool *sqlitex.Pool
}

// Conn returns one connection from connection pool
// NOTE: make sure to call done function to put the connection back to the pool
func (db *Database) Conn(ctx context.Context) (conn *sqlite.Conn, done func(), err error) {
	conn = db.pool.Get(ctx)
	if conn == nil {
		return nil, nil, context.Canceled
	}

	return conn, func() { db.pool.Put(conn) }, nil
}

type Options struct {
	StringConn string
	PoolSize   int
}

// New creates a sqlite database
func New(opt Options) (*Database, error) {
	pool, err := sqlitex.Open(opt.StringConn, 0, opt.PoolSize)
	if err != nil {
		return nil, err
	}

	// the following loop makes sure that all pool connections have
	// forgen_key enabled by default
	for i := 0; i < opt.PoolSize; i++ {
		conn := pool.Get(context.Background())
		err := sqlitex.Exec(conn, `PRAGMA foreign_keys = ON;`, nil)
		if err != nil {
			return nil, err
		}
		pool.Put(conn)
	}

	return &Database{
		pool: pool,
	}, nil
}

func RunScript(conn *sqlite.Conn, sql string) error {
	return sqlitex.ExecScript(conn, strings.TrimSpace(sql))
}
