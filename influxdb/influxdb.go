package influxdb

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/vgmdj/utils/logger"

	"github.com/influxdata/influxdb/client"
)

//InfluxClient the struct of influxDB client
type InfluxClient struct {
	cli *client.Client
}

var (
	c    *InfluxClient
	once sync.Once
)

//NewClient create a global client with server addr including schema, host, and ip
func NewClient(server string) *InfluxClient {

	once.Do(func() {
		host, err := url.Parse(fmt.Sprintf("http://%s:%d", "localhost", 8086))
		if err != nil {
			logger.Error(err.Error())
		}
		con, err := client.NewClient(client.Config{URL: *host})
		if err != nil {
			logger.Error(err.Error())
		}

		c = &InfluxClient{
			cli: con,
		}

	})

	return c
}
