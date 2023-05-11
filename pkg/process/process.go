package process

import (
	"github.com/notivenl/uptime-kubernetes/pkg/config"
	v1 "k8s.io/api/networking/v1"
)

// Data is the struct that is passed between the inputs and outputs
// It contains all the HTTP data that is needed to perform the health check
type Data struct {
	Name     string
	Protocol string
	Method   string
	Url      string
	Body     string
	Headers  map[string]string

	HealthStatus int
}

//go:generate moq -pkg mock -out ../../mocks/process/input.go . Input

// Input is the interface that all inputs need to implement
// Inputs are only responsible for getting the data from the ingress and sending it to the data channel
type Input interface {
	Init(*config.Config) error
	Process(v1.Ingress, chan Data)
}

//go:generate moq -pkg mock -out ../../mocks/process/output.go . Output

// Output is the interface that all outputs need to implement
// Outputs are only responsible for performing an action using the data from the data channel if any is available
type Output interface {
	Init(*config.Config) error
	Process(Data)
}

var inputs map[string]Input
var outputs map[string]Output

// RegisterInput registers an input
// It panics if the input is already registered
func RegisterInput(name string, input Input) {
	if _, ok := inputs[name]; ok {
		panic("input already registered: " + name)
	}
	inputs[name] = input
}

// RegisterOutput registers an output
// It panics if the output is already registered
func RegisterOutput(name string, output Output) {
	if _, ok := outputs[name]; ok {
		panic("output already registered: " + name)
	}
	outputs[name] = output
}

// GetInput returns an input if it exists with the given name or nil
func GetInput(name string) Input {
	if input, ok := inputs[name]; ok {
		return input
	}
	return nil
}

// GetOutput returns an output if it exists with the given name or nil
func GetOutput(name string) Output {
	if output, ok := outputs[name]; ok {
		return output
	}
	return nil
}

func init() {
	inputs = make(map[string]Input)
	outputs = make(map[string]Output)
}
