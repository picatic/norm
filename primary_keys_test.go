package norm

import (
	"github.com/picatic/norm/field"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPrimaryKeyer(t *testing.T) {
	Convey("PrimaryKeyer", t, func() {

		Convey("SinglePrimaryKey", func() {
			pk := NewSinglePrimaryKey(field.Name("Id"))
			So(pk.Fields(), ShouldResemble, field.Names{field.Name("Id")})
		})

		Convey("MultiplePrimaryKey", func() {
			pk := NewMultiplePrimaryKey(field.Names{"Id", "Org"})
			So(pk.Fields(), ShouldResemble, field.Names{"Id", "Org"})
		})
	})
}
