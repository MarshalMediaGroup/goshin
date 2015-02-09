package checks


import (
	"github.com/MarshalMediaGroup/goshin"
	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
	"fmt"
)


type ServerStatus struct {
	db *mgo.Database
}

func NewServerStatus(db *mgo.Database) *ServerStatus {
	return &ServerStatus{db:db}
}


func (m *ServerStatus) Collect(queue chan *goshin.Metric) {
	var serverStatus = bson.M{};
	if err := m.db.Run("serverStatus", &serverStatus); err!=nil{
		fmt.Println(err)
		return
	}
	m.collectFromBSON(&serverStatus, "serverStatus", queue)
}

func (m *ServerStatus) collectFromBSON(source *bson.M, prefix string, queue chan *goshin.Metric) {
	for key, value := range *source{
		name := fmt.Sprintf("%s.%s", prefix, key)
		switch value := value.(type) {
		default:
		case int, float64:
			queue <- m.buildMetric(name, value)
		case int64:
			queue <- m.buildMetric(name, int(value))
		case bson.M:
			m.collectFromBSON(&value, name, queue)
		}
	}
}
func (m *ServerStatus) buildMetric(service string, value interface {}) *goshin.Metric{
	metric := goshin.NewMetric()
	metric.Service = service
	metric.Value = value
	return metric
}

