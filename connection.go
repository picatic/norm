package norm

import (
	"database/sql"
	"github.com/gocraft/dbr"
)

type Connection struct {
	*dbr.Connection
	database string
}

func NewConnection(db *sql.DB, database string, log dbr.EventReceiver) *Connection {
	return &Connection{Connection: dbr.NewConnection(db, log), database: database}
}

func (c Connection) Database() string {
	return c.database
}

func (c Connection) NewSession(log dbr.EventReceiver) *Session {
	return &Session{Session: c.Connection.NewSession(log), connection: &c}
}

type Session struct {
	*dbr.Session
	connection *Connection
}

func (s Session) Connection() *Connection {
	return s.connection
}
