// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/notivenl/uptime-kubernetes/pkg/config"
	"github.com/notivenl/uptime-kubernetes/pkg/process"
	"sync"
)

// Ensure, that OutputMock does implement process.Output.
// If this is not the case, regenerate this file with moq.
var _ process.Output = &OutputMock{}

// OutputMock is a mock implementation of process.Output.
//
// 	func TestSomethingThatUsesOutput(t *testing.T) {
//
// 		// make and configure a mocked process.Output
// 		mockedOutput := &OutputMock{
// 			InitFunc: func(configMoqParam *config.Config) error {
// 				panic("mock out the Init method")
// 			},
// 			ProcessFunc: func(data process.Data)  {
// 				panic("mock out the Process method")
// 			},
// 		}
//
// 		// use mockedOutput in code that requires process.Output
// 		// and then make assertions.
//
// 	}
type OutputMock struct {
	// InitFunc mocks the Init method.
	InitFunc func(configMoqParam *config.Config) error

	// ProcessFunc mocks the Process method.
	ProcessFunc func(data process.Data)

	// calls tracks calls to the methods.
	calls struct {
		// Init holds details about calls to the Init method.
		Init []struct {
			// ConfigMoqParam is the configMoqParam argument value.
			ConfigMoqParam *config.Config
		}
		// Process holds details about calls to the Process method.
		Process []struct {
			// Data is the data argument value.
			Data process.Data
		}
	}
	lockInit    sync.RWMutex
	lockProcess sync.RWMutex
}

// Init calls InitFunc.
func (mock *OutputMock) Init(configMoqParam *config.Config) error {
	if mock.InitFunc == nil {
		panic("OutputMock.InitFunc: method is nil but Output.Init was just called")
	}
	callInfo := struct {
		ConfigMoqParam *config.Config
	}{
		ConfigMoqParam: configMoqParam,
	}
	mock.lockInit.Lock()
	mock.calls.Init = append(mock.calls.Init, callInfo)
	mock.lockInit.Unlock()
	return mock.InitFunc(configMoqParam)
}

// InitCalls gets all the calls that were made to Init.
// Check the length with:
//     len(mockedOutput.InitCalls())
func (mock *OutputMock) InitCalls() []struct {
	ConfigMoqParam *config.Config
} {
	var calls []struct {
		ConfigMoqParam *config.Config
	}
	mock.lockInit.RLock()
	calls = mock.calls.Init
	mock.lockInit.RUnlock()
	return calls
}

// Process calls ProcessFunc.
func (mock *OutputMock) Process(data process.Data) {
	if mock.ProcessFunc == nil {
		panic("OutputMock.ProcessFunc: method is nil but Output.Process was just called")
	}
	callInfo := struct {
		Data process.Data
	}{
		Data: data,
	}
	mock.lockProcess.Lock()
	mock.calls.Process = append(mock.calls.Process, callInfo)
	mock.lockProcess.Unlock()
	mock.ProcessFunc(data)
}

// ProcessCalls gets all the calls that were made to Process.
// Check the length with:
//     len(mockedOutput.ProcessCalls())
func (mock *OutputMock) ProcessCalls() []struct {
	Data process.Data
} {
	var calls []struct {
		Data process.Data
	}
	mock.lockProcess.RLock()
	calls = mock.calls.Process
	mock.lockProcess.RUnlock()
	return calls
}
