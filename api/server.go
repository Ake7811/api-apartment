package api

import (
	"fmt"

	
)

func Run() {
	config.Load();
	fmt.Printf("running... at port %d", config.PORT)
}