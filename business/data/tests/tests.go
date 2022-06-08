// Package tests contains supporting code for running tests
package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/yakushou730/ardanlabs-ultimate-serice-v3/foundation/logger"

	"github.com/yakushou730/ardanlabs-ultimate-serice-v3/business/data/schema"

	"github.com/yakushou730/ardanlabs-ultimate-serice-v3/business/sys/database"

	"github.com/yakushou730/ardanlabs-ultimate-serice-v3/foundation/docker"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Success and failure markers
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// DBContainer provides configuration for a container to run
type DBContainer struct {
	Image string
	Port  string
	Args  []string
}

// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewUnit(t *testing.T, dbc DBContainer) (*zap.SugaredLogger, *sqlx.DB, func()) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	c := docker.StartContainer(t, dbc.Image, dbc.Port, dbc.Args...)

	db, err := database.Open(
		database.Config{
			User:         "root",
			Password:     "secret",
			Host:         "local-postgresql.default.svc:5432", // "localhost:30001",
			Name:         "sales_dev",
			MaxIdleConns: 0,
			MaxOpenConns: 0,
			DisableTLS:   true,
		})
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("waiting for database to be ready...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := schema.Migrate(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		docker.StopContainer(t, c.ID)
		t.Fatalf("Migrating error: %s", err)
	}

	if err := schema.Seed(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		docker.StopContainer(t, c.ID)
		t.Fatalf("Seeding error: %s", err)
	}

	log, err := logger.New("TEST")
	if err != nil {
		t.Fatalf("logger error: %s", err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
		docker.StopContainer(t, c.ID)

		log.Sync()

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old
		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}

	return log, db, teardown
}
