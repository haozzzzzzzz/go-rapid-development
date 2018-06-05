package printutil

import (
	"fmt"
)

func PrintObject(obj interface{}) {
	fmt.Printf("%#v\n", obj)
}
