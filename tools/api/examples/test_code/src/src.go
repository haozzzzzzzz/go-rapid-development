package src

import (
	"fmt"
	// +pre if(USE_IMPORT) then
	_ "fmt"
	// +pre else
	_ "log"
	// +pre end
)

func SaySomething() {
	fmt.Println("hello")
}
