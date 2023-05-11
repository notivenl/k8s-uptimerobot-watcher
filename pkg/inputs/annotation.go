package inputs

import (
	"strconv"
	"strings"

	"github.com/notivenl/uptime-kubernetes/pkg/config"
	"github.com/notivenl/uptime-kubernetes/pkg/process"
	v1 "k8s.io/api/networking/v1"
)

// Annotation is an input type that reads the annotations from the ingress
// and puts the data from them in the Data struct
type Annotation struct{}

func init() {
	process.RegisterInput("annotation", &Annotation{})
}

// Init is called when the input is started
func (a *Annotation) Init(config *config.Config) error {
	return nil
}

// Process gets the data from the ingress and puts it in the data channel
func (a *Annotation) Process(ingress v1.Ingress, dataChan chan process.Data) {
	for _, value := range ingress.Spec.Rules {
		data := &process.Data{
			Headers: make(map[string]string),
		}

		for key, value := range ingress.Annotations {
			err := a.annotationKey(data, key, value)
			if err != nil {
				panic(err)
			}
		}

		data.Url = value.Host + data.Url

		dataChan <- *data
	}
}

const NOTIVE_PREFIX = "io.notive.health/"

// annotationKey puts the data from the annotation in the Data struct in its respective fields
// example annotations:
// io.notive.health/name: <name>
// io.notive.health/protocol: https
// io.notive.health/method: GET, POST, PUT, DELETE
// io.notive.health/check: /health
// io.notive.health/body: {"status": "ok"}
// io.notive.health/status: 200
// io.notive.health/header-<headername>: <headervalue>
func (a *Annotation) annotationKey(data *process.Data, key string, value string) error {
	switch key {
	case NOTIVE_PREFIX + "name":
		data.Name = value
	case NOTIVE_PREFIX + "protocol":
		data.Protocol = value
	case NOTIVE_PREFIX + "method":
		data.Method = value
	case NOTIVE_PREFIX + "check":
		data.Url = value
	case NOTIVE_PREFIX + "body":
		data.Body = value
	case NOTIVE_PREFIX + "status":
		status, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		data.HealthStatus = status
	default:
		if strings.HasPrefix(key, NOTIVE_PREFIX+"header-") {
			key = strings.TrimPrefix(key, NOTIVE_PREFIX+"header-")
			data.Headers[key] = value
		}
	}
	return nil
}
