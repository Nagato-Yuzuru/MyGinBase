package main

import (
	"GinBase/pkg/config"
	"fmt"
)

func main() {
	loader := config.NewViperLoader("config", "yaml", []string{"./config"}, "", nil)
	cfg, err := loader.Load()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.App)
}
