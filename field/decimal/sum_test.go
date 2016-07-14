package decimal

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSum(t *testing.T) {
	Convey("Sum", t, func() {
		list := []interface{}{
			Dec{Number: 123, Prec: 2},
			Dec{Number: 456, Prec: 2},
			Dec{Number: 789, Prec: 2},
		}

		sum := Sum(list)
		So(sum.String(), ShouldEqual, "13.68")
	})
}
