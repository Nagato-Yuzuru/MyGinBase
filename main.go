package main

import (
	"GinBase/inject"
	"fmt"
)

func main() {
	config := inject.InitializeConfig()
	fmt.Println(config)
}
