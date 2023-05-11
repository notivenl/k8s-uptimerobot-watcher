package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/notivenl/uptime-kubernetes/pkg/config"
	_ "github.com/notivenl/uptime-kubernetes/pkg/inputs"
	_ "github.com/notivenl/uptime-kubernetes/pkg/outputs"
	"github.com/notivenl/uptime-kubernetes/pkg/ports/kubernetes"
	"github.com/notivenl/uptime-kubernetes/pkg/process"
	v1i "k8s.io/api/networking/v1"
)

// UptimeClient is the main client for the uptime-kubernetes application.
// It is responsible for fetching the ingresses from the Kubernetes API and sending them to the input.
// It is also responsible for receiving the data from the output and sending it to the Kubernetes API.
type UptimeClient struct {
	Context context.Context

	Config *config.Config
	Client kubernetes.Kubernetes

	ingressChan chan v1i.Ingress
	dataChan    chan process.Data

	input  process.Input
	output process.Output

	Timeout time.Duration

	tickMutex sync.Mutex
	stop      chan bool
}

// NewUptimeClient creates a new UptimeClient.
func NewUptimeClient(applicationConfig *config.Config, rest kubernetes.Rest) *UptimeClient {

	client, err := rest.GetInternalClient()
	if err != nil {
		panic(err)
	}

	input := process.GetInput(applicationConfig.Input)
	if input == nil {
		panic(fmt.Sprintf("Cannot find input with the name of '%s'", applicationConfig.Input))
	}
	output := process.GetOutput(applicationConfig.Output)
	if output == nil {
		panic(fmt.Sprintf("Cannot find output with the name of '%s'", applicationConfig.Output))
	}

	input.Init(applicationConfig)
	output.Init(applicationConfig)

	return &UptimeClient{
		Context: context.Background(),

		Config: applicationConfig,

		Client: client,

		ingressChan: make(chan v1i.Ingress),
		dataChan:    make(chan process.Data),

		input:  input,
		output: output,

		Timeout: time.Duration(applicationConfig.TickSpeed) * time.Second,
		stop:    make(chan bool),
	}
}

// Start starts the main loop of the UptimeClient using a ticker, the ticker speed can be configured in the config.
// if theres data in the ingress channel or in the data channel it will start and input or output process in a seperate goroutine.
func (c *UptimeClient) Start() {
	ticker := time.NewTicker(c.Timeout)

	fmt.Println("Starting UptimeClient...")

	for {
		select {
		case <-ticker.C:
			fmt.Println("Checking ingresses...")
			go c.Tick() // needs to be in seperate goroutine to prevent blocking

		case ingress := <-c.ingressChan:
			fmt.Println("Processing ingress...")
			go c.input.Process(ingress, c.dataChan)

		case data := <-c.dataChan:
			fmt.Println("Processing data...")
			go c.output.Process(data)

		case <-c.stop:
			ticker.Stop()
			return
		}
	}
}

// Stop stops the main loop of the UptimeClient.
func (c *UptimeClient) Stop() {
	c.stop <- true
}

// IngressChan requests the ingresses from the Kubernetes API and feeds them into the input channel.
func (c *UptimeClient) Tick() {
	if c.tickMutex.TryLock() { // prevents multiple ticks from running at the same time
		ingresses, err := c.Client.GetIngresses(c.Context)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, ingress := range ingresses {
			c.ingressChan <- ingress
		}
		c.tickMutex.Unlock()
	}
}
