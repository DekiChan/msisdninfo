package main

import (
	"fmt"

	"github.com/dekichan/msisdninfo/serve"
)

func main() {
	fmt.Println("msisdninfo started")

	serve.Serve()
}
