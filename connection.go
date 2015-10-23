package norm

import (
	"database/sql"
	"github.com/gocraft/dbr"
)

// Connection initialized with a database
type Connection interface {
	// dbr.Connection
	NewSession(log dbr.EventReceiver) Session

	Database() string
	ValidatorCache() ValidatorCache
}

type connection struct {
	*dbr.Connection
	database       string
	validatorCache ValidatorCache
}

// NewConnection return a Connection as configured
func NewConnection(db *sql.DB, database string, log dbr.EventReceiver) Connection {
	return &connection{Connection: dbr.NewConnection(db, log), database: database, validatorCache: make(ValidatorCache, 0)}
}

// Database returns name of database
func (c connection) Database() string {
	return c.database
}

// ValidatorCache returns ValidatorCache
func (c connection) ValidatorCache() ValidatorCache {
	return c.validatorCache
}

// NewSession Create a new Session with the Connection
func (c connection) NewSession(log dbr.EventReceiver) Session {
	return &session{Session: c.Connection.NewSession(log), connection: &c}
}

// Session return a Session to work with
type Session interface {
	// dbr.Session functions
	Begin() (*dbr.Tx, error)
	DeleteFrom(from string) *dbr.DeleteBuilder
	InsertInto(into string) *dbr.InsertBuilder
	Select(cols ...string) *dbr.SelectBuilder
	SelectBySql(sql string, args ...interface{}) *dbr.SelectBuilder
	Update(table string) *dbr.UpdateBuilder
	UpdateBySql(sql string, args ...interface{}) *dbr.UpdateBuilder

	Connection() Connection
}

type session struct {
	*dbr.Session
	connection Connection
}

// Connection returns the connection used to create the session
func (s session) Connection() Connection {
	return s.connection
}

type Tx struct {
	*dbr.Tx
	connection Connection
}

func (t Tx) Connection() Connection {
	return t.connection
}

func Begin(s Session) (*Tx, error) {
	tx, err := s.Begin()
	return &Tx{tx, s.Connection()}, err
}
