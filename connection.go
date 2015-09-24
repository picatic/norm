package norm

import (
	"database/sql"
	"github.com/gocraft/dbr"
)

// Connection initialized with a database
type Connection struct {
	*dbr.Connection
	database string
}

// NewConnection return a Connection as configured
func NewConnection(db *sql.DB, database string, log dbr.EventReceiver) *Connection {
	return &Connection{Connection: dbr.NewConnection(db, log), database: database}
}

// Database returns name of database
func (c Connection) Database() string {
	return c.database
}

// NewSession Create a new Session with the Connection
func (c Connection) NewSession(log dbr.EventReceiver) *Session {
	return &Session{Session: c.Connection.NewSession(log), connection: &c}
}

// Session return a Session to work with
type Session struct {
	*dbr.Session
	connection *Connection
}

// Connection returns the connection used to create the session
func (s Session) Connection() *Connection {
	return s.connection
}
