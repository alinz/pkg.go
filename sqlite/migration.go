package sqlite

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type Migration interface {
	Name() string
	Skip() bool
	Up(conn *sqlite.Conn) error
	Down(conn *sqlite.Conn) error
}

type MigrationRunner struct {
	db         *Database
	migrations []Migration
}

func (m *MigrationRunner) perform(ctx context.Context, direction string) (err error) {
	conn, closeConn, err := m.db.Conn(ctx)
	if err != nil {
		return nil
	}
	defer closeConn()
	defer sqlitex.Save(conn)(&err)

	log.
		Debug().
		Int("migrations", len(m.migrations)).
		Msg("Start applying migrations")

	switch direction {
	case "up":
		// Sorts the Migration based on their names in ascending order
		// because we need to go up
		sort.SliceStable(m.migrations, func(i, j int) bool {
			return m.migrations[i].Name() < m.migrations[j].Name()
		})

		for _, migration := range m.migrations {
			if migration.Skip() {
				log.
					Debug().
					Str("name", migration.Name()).
					Msg("skip by force")
				continue
			}

			applied, err := m.has(conn, migration.Name())
			if err != nil {
				log.
					Error().
					Err(err).
					Str("name", migration.Name()).
					Str("direction", "up").
					Msg("something went wrong during migration")
				return err
			}

			// if migration already applied, then skip
			// the migration up
			if applied {
				log.
					Debug().
					Str("name", migration.Name()).
					Str("direction", "up").
					Msg("skip as it's has already been applied")
				continue
			}

			err = migration.Up(conn)
			if err != nil {
				log.
					Error().
					Err(err).
					Str("name", migration.Name()).
					Str("direction", "up").
					Msg("something went wrong during migration")
				return err
			}

			err = m.add(conn, migration.Name())
			if err != nil {
				log.
					Error().
					Err(err).
					Str("name", migration.Name()).
					Str("direction", "up").
					Msg("something went wrong during migration")
				return err
			}
		}

	case "down":
		// Sorts the Migration based on their names in descending order
		// because we need to go down
		sort.SliceStable(m.migrations, func(i, j int) bool {
			return m.migrations[i].Name() > m.migrations[j].Name()
		})

		for _, migration := range m.migrations {
			if migration.Skip() {
				log.
					Debug().
					Str("name", migration.Name()).
					Msg("skip by force")
				continue
			}

			applied, err := m.has(conn, migration.Name())
			if err != nil {
				log.
					Error().
					Err(err).
					Str("name", migration.Name()).
					Str("direction", "down").
					Msg("something went wrong during migration")
				return err
			}

			// if migration hasn't happened yet, then skip
			// the migration down
			if !applied {
				log.
					Debug().
					Str("name", migration.Name()).
					Str("direction", "down").
					Msg("skip as it's hasn't been applied yet")
				continue
			}

			err = migration.Down(conn)
			if err != nil {
				log.
					Error().
					Err(err).
					Str("name", migration.Name()).
					Str("direction", "down").
					Msg("something went wrong during migration")
				return err
			}

			err = m.remove(conn, migration.Name())
			if err != nil {
				log.
					Error().
					Err(err).
					Str("name", migration.Name()).
					Str("direction", "down").
					Msg("something went wrong during migration")
				return err
			}
		}

	default:
		return fmt.Errorf("wrong direction, %s", direction)
	}

	return nil
}

func (m *MigrationRunner) remove(conn *sqlite.Conn, name string) error {
	stmt, err := conn.Prepare(`DELETE FROM migrations WHERE name=$name;`)
	if err != nil {
		return err
	}

	defer stmt.Finalize()

	stmt.SetText("$name", name)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	return nil
}

func (m *MigrationRunner) add(conn *sqlite.Conn, name string) error {
	stmt, err := conn.Prepare(`INSERT INTO migrations (name) VALUES ($name) ON CONFLICT(name) DO UPDATE SET name=$u_name;`)
	if err != nil {
		return err
	}

	defer stmt.Finalize()

	stmt.SetText("$name", name)
	stmt.SetText("$u_name", name)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	return nil
}

func (m *MigrationRunner) has(conn *sqlite.Conn, name string) (bool, error) {
	stmt, err := conn.Prepare(`SELECT name from migrations WHERE name=$name;`)
	if err != nil {
		return false, err
	}

	defer stmt.Finalize()

	stmt.SetText("$name", name)

	rowReturned, err := stmt.Step()
	if err != nil {
		return false, err
	}

	return rowReturned, nil
}

// Up goes through all migration in upward direction, if any of the
// migration step fails, it returns an error. The whole up is protected
// by transaction. So all need to be success in order for up to be success
func (m *MigrationRunner) Up(ctx context.Context) error {
	return m.perform(ctx, "up")
}

// Down goes through all migration in downward direction, if any of the
// migration step fails, it returns an error. The whole down is protected
// by transaction. So all need to be success in order for down to be success
func (m *MigrationRunner) Down(ctx context.Context) error {
	return m.perform(ctx, "down")
}

// NewMigrationRunner creates a runner that facilitates up and down migration
func NewMigrationRunner(db *Database, migrations ...Migration) (*MigrationRunner, error) {
	conn, closeConn, err := db.Conn(context.Background())
	if err != nil {
		log.
			Error().
			Err(err).
			Msg("failed to get conn for creating migrations table")
		return nil, err
	}
	defer closeConn()

	err = sqlitex.ExecScript(conn, strings.TrimSpace(`
		CREATE TABLE IF NOT EXISTS migrations(
			name TEXT NOT NULL,
			PRIMARY KEY (name)
		);
	`))
	if err != nil {
		log.
			Error().
			Err(err).
			Msg("failed to create migrations table")
		return nil, err
	}

	return &MigrationRunner{
		db:         db,
		migrations: migrations,
	}, nil
}
