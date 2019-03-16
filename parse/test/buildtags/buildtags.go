// +build x,y z
// +build genny

package buildtags

import (
	"fmt"

	"github.com/mauricelam/genny/generic"
)

type _t_ generic.Type

func _t_Print(t _t_) {
	fmt.Println(t)
}
