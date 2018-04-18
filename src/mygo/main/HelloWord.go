package main

import (
	"fmt"
	"os"
)

func main() {
	var s ,sep string
	for i := 0; i < os.Args[i]; i++ {
		s += sep + os.Args[i];
		sep = "";
	}
	fmt.Println(s);
}

