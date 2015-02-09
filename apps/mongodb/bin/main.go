package main

import (
	"github.com/MarshalMediaGroup/goshin/apps/mongodb"
)

func main() {
	app := mongodb.New()
	app.Start()
}
