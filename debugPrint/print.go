package debugPrint

import (
	"fmt"
)

var Println func(a ...any)

func PrintNoop(a ...any) {}

func PrintFmt(a ...any) {
	fmt.Print("DEBUG: ")
	fmt.Println(a...)
}
