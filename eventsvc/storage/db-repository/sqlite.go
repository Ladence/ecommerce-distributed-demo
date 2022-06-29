package db_repository

import (
	"context"
	"database/sql"
	"github.com/Ladence/ecommerce-distributed-demo/eventstore"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"strings"
)

type sqliteRepository struct {
	db *sql.DB
}

// todo: отрефакорить асап и инжектить в NewSqliteRepository *sql.DB
func NewDbRepository(db *sql.DB) *sqliteRepository {
	return &sqliteRepository{
		db: db,
	}
}

func NewSqliteRepository(connectionString string) (*sqliteRepository, error) {
	dbEvents, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return nil, err
	}
	if err = setupSchema(dbEvents); err != nil {
		return nil, err
	}
	return &sqliteRepository{
		db: dbEvents,
	}, nil
}

func setupSchema(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}
	fi, err := (&file.File{}).Open("./migrations/events/")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("file", fi, "sqlite3", driver)
	if err != nil {
		return err
	}
	err = m.Steps(1)
	if err != nil {
		return err
	}
	return nil
}

func (d *sqliteRepository) CreateEvent(ctx context.Context, event *eventstore.Event) error {
	var e error
	query := "INSERT INTO events(id, eventtype, aggregateid, aggregatetype, eventdata, stream) VALUES($1, $2, $3, $4, $5, $6);"
	if event.EventData == "" {
		query = strings.Replace(query, "$5", "NULL", 1)
		query = strings.Replace(query, "$6", "$5", 1)
		_, e = d.db.ExecContext(ctx, query, event.EventId, event.EventType, event.AggregateId, event.AggregateType, event.Stream)
	} else {
		_, e = d.db.ExecContext(ctx, query, event.EventId, event.EventType, event.AggregateId, event.AggregateType, event.EventData, event.Stream)
	}
	return e
}

func (d *sqliteRepository) GetEvents(ctx context.Context, filter *eventstore.GetEventsRequest) ([]*eventstore.Event, error) {
	var rows *sql.Rows
	var e error
	var query string

	if filter.EventId == "" && filter.AggregateId == "" {
		query = "SELECT id, eventtype, aggregateid, aggregatetype, eventdata FROM events;"
		rows, e = d.db.QueryContext(ctx, query)
	} else if filter.EventId != "" && filter.AggregateId == "" {
		query = "SELECT id, eventtype, aggregateid, aggregatetype, eventdata FROM events WHERE id=$1;"
		rows, e = d.db.QueryContext(ctx, query, filter.EventId)
	} else if filter.EventId == "" && filter.AggregateId != "" {
		query = "SELECT id, eventtype, aggregateid, aggregatetype, eventdata FROM events WHERE aggregateid=$1;"
		rows, e = d.db.QueryContext(ctx, query, filter.AggregateId)
	} else if filter.EventId != "" && filter.AggregateId != "" {
		query = "SELECT id, eventtype, aggregateid, aggregatetype, eventdata FROM events WHERE id=$1 AND aggregateid=$2;"
		rows, e = d.db.QueryContext(ctx, query, filter.EventId, filter.AggregateId)
	}
	if e != nil {
		return nil, e
	}

	events := make([]*eventstore.Event, 0)
	for rows.Next() {
		event, e := scanEvent(rows)
		if e != nil {
			return events, e
		}
		events = append(events, event)
	}
	return events, e
}

func scanEvent(rows *sql.Rows) (*eventstore.Event, error) {
	event := &eventstore.Event{}
	if err := rows.Scan(&event.EventId, &event.EventType, &event.AggregateId, &event.AggregateType, &event.EventData); err != nil {
		return nil, err
	}
	return event, nil
}
