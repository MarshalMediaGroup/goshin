package mongodb

import (
	"flag"
	"fmt"
	"github.com/MarshalMediaGroup/goshin"
	"github.com/MarshalMediaGroup/goshin/apps/mongodb/checks"
	mgo "gopkg.in/mgo.v2"
)

type MongoDb struct {
	*goshin.Goshin

	db *mgo.Database
}

func New() *MongoDb {
	app := &MongoDb{}
	app.Goshin = goshin.NewGoshin()
	app.Configure()

	app.AddCheck("serverStatus", checks.NewServerStatus(app.db))
	return app
}

func (app *MongoDb) Configure() {
	var (
		mongoHostPtr = flag.String("mongo-host", "localhost", "Mongo hostname")
		mongoPortPtr = flag.Int("mongo-port", 27017, "Mongo port")
		mongoDBPtr   = flag.String("mongo-db", "local", "Mongo database")
		checksPtr    = flag.String("checks", "serverStatus", "A list of checks to run")
	)
	app.Goshin.Configure()
	app.ExtractEnabledChecks(*checksPtr)

	url := fmt.Sprintf("%s:%d", *mongoHostPtr, *mongoPortPtr)
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	app.db = session.DB(*mongoDBPtr)
}
