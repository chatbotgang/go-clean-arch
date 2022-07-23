package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var testPostgresDB *sqlx.DB

func getTestPostgresDB() *sqlx.DB {
	return testPostgresDB
}

// migrationSourcePath is a relative path to the collection storing all migration scripts
const migrationSourcePath = "file://../../../../migrations"
const testPostgresName = "repo_test"

func buildTestPostgresDB() (*sqlx.DB, func(), error) {
	cb, pgdsn, err := startPostgresContainer()
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to start postgres container")
	}

	// migrate postgres
	m, err := migrate.New(migrationSourcePath, pgdsn)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to start migration process")
	}

	err = m.Up()
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to migrate postgres")
	}

	// connect to postgres
	db, err := sqlx.Open("postgres", pgdsn)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to connect postgres")
	}
	return db, cb, nil
}

func startPostgresContainer() (func(), string, error) {
	// new a docker pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, "", fmt.Errorf("could not connect to docker: %s", err)
	}
	// start a PG container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.1",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", testPostgresName),
			fmt.Sprintf("POSTGRES_USER=%s", testPostgresName),
			fmt.Sprintf("POSTGRES_DB=%s", testPostgresName),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, "", fmt.Errorf("could not start postgres: %s", err)
	}

	// Get host and port(random) info from the postgres container
	hostAndPort := resource.GetHostPort("5432/tcp")
	LocalhostPostgresDSN := fmt.Sprintf(
		"postgresql://repo_test:repo_test@%s/repo_test?sslmode=disable",
		hostAndPort,
	)

	// build a call back function to destroy the docker pool
	cb := func() {
		if err := pool.Purge(resource); err != nil {
			log.Printf("Could not purge resource: %s", err)
		}
	}
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err := sql.Open("postgres", LocalhostPostgresDSN)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		cb()
		return nil, "", fmt.Errorf("could not connect to postgres: %s", err)
	}

	return cb, LocalhostPostgresDSN, nil
}

func initRepository(t *testing.T, db *sqlx.DB, files ...string) (repo *PostgresRepository) {
	// Truncate existing records in all tables
	truncateAllData(t, db)

	// Setup DB again
	loader, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.Location(time.UTC),
		// Load predefined data
		testfixtures.Files(files...),
	)
	require.NoError(t, err)

	err = loader.Load()
	require.NoError(t, err)

	return NewPostgresRepository(context.Background(), db)
}

func truncateAllData(t *testing.T, db *sqlx.DB) {
	template := `CREATE OR REPLACE FUNCTION truncate_all_tables() RETURNS void AS $$
				DECLARE
					statements CURSOR FOR
						SELECT tablename FROM %s.pg_catalog.pg_tables
						WHERE schemaname = 'public' AND tablename !='schema_migrations';
				BEGIN
					FOR stmt IN statements LOOP
						EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.tablename) || ' RESTART IDENTITY CASCADE;';
					END LOOP;
				END;
				$$ LANGUAGE plpgsql;
				SELECT truncate_all_tables();`

	// We add test postgres db name to avoid mis-using the script in production.
	script := fmt.Sprintf(template, testPostgresName)
	_, err := db.Exec(script)
	require.NoError(t, err)
}

// nolint
func TestMain(m *testing.M) {
	// To avoid violating table constraints
	_ = faker.SetRandomStringLength(10)
	_ = faker.SetRandomMapAndSliceMinSize(10)

	db, closeDB, err := buildTestPostgresDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer closeDB()
	testPostgresDB = db
	m.Run()
}
