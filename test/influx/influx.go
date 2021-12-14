package influx

import (
	// "Hybrid_Cluster/omcplog"
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
)

type Influx struct {
	inClient client.Client
}

func NewInflux(INFLUX_IP, INFLUX_PORT, username, password string) *Influx {
	// omcplog.V(4).Info("Func NewInflux Called")
	inf := &Influx{
		inClient: InfluxDBClient(INFLUX_IP, INFLUX_PORT, username, password),
	}
	return inf
}

func InfluxDBClient(INFLUX_IP, INFLUX_PORT, username, password string) client.Client {
	// omcplog.V(4).Info("Func InfluxDBClient Called")
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + INFLUX_IP + ":" + INFLUX_PORT,
		Username: username,
		Password: password,
	})
	if err != nil {
		// omcplog.V(0).Info("Error: ", err)
		panic(err)
	}
	return c
}

func (in *Influx) GetPodData(podNum int, ns, clusterName string) []client.Result {
	// omcplog.V(4).Info("Func GetPodData Called")
	// omcplog.V(5).Info("timeStart=", timeStart, ", timeEnd=", timeEnd)
	pod_Num := fmt.Sprint(podNum)
	q := client.NewQuery("SELECT * FROM Pods WHERE cluster = '"+clusterName+"' ORDER BY DESC LIMIT "+pod_Num+" ", "Metrics", "")

	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {
		return response.Results
	}
	return nil
}

func (in *Influx) GetNodeData(clusterName string) []client.Result {
	// omcplog.V(4).Info("Func GetNodeData Called")
	// omcplog.V(5).Info("timeStart=", timeStart, ", timeEnd=", timeEnd)

	q := client.NewQuery("SELECT * FROM Nodes WHERE cluster = '"+clusterName+"' ORDER BY DESC LIMIT 1", "Metrics", "")

	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {
		return response.Results
	}
	return nil

}
