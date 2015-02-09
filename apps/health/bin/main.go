package main

import (
	"github.com/MarshalMediaGroup/goshin/apps/health"
)

func main() {
	app := health.New()
	app.Start()
}
