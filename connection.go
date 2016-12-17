package norm

import (
	"database/sql"
	"errors"
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"
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
	conn := &dbr.Connection{}
	conn.DB = db
	conn.Dialect = dialect.MySQL
	conn.EventReceiver = log
	return &connection{Connection: conn, database: database, validatorCache: make(ValidatorCache, 0)}
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
	Begin() (Tx, error)
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

// Begin returns a norm Tx which has wrapped a dbr.Tx
// A real database connection has been aquired and is held by the enclosed sql.Tx instance
func (s session) Begin() (Tx, error) {
	dbrTx, err := s.Session.Begin()
	return &tx{dbrTx, s.Connection()}, err
}

// Tx embeds dbr.Tx and norm Session
type Tx interface {
	Session
	Commit() error
	Rollback() error
	RollbackUnlessCommitted()
}

// tx implements Tx interface
type tx struct {
	*dbr.Tx
	connection Connection
}

// Connection returns norm Connection
func (t tx) Connection() Connection {
	return t.connection
}

func (t tx) Begin() (Tx, error) {
	return nil, errors.New("Support for nested transactions not implemented")
}
