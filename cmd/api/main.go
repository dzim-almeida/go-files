package main

import "github.com/dzim-almeida/go-files/internal/router"

func main() {
	r := router.NewRouter()

	r.Run(":8000")
}