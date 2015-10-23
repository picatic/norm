package norm

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConnection(t *testing.T) {
	Convey("Connection", t, func() {
		db, _ := sql.Open("", "")
		conn := NewConnection(db, "mocked", nil)
		Convey("Has Database()", func() {
			So(conn.Database(), ShouldEqual, "mocked")
		})
	})

	Convey("Session", t, func() {
		db, _ := sql.Open("", "")
		conn := NewConnection(db, "mocked", nil)
		sess := conn.NewSession(nil)

		Convey("Has reference to Connection", func() {
			So(sess.Connection(), ShouldResemble, conn)
		})
	})
}

var tx Session = Tx{} //ensure Tx implements Session
