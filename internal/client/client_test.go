package client_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/notivenl/uptime-kubernetes/internal/client"
	mockKubernetes "github.com/notivenl/uptime-kubernetes/mocks/kubernetes"
	mock "github.com/notivenl/uptime-kubernetes/mocks/process"
	"github.com/notivenl/uptime-kubernetes/pkg/config"
	"github.com/notivenl/uptime-kubernetes/pkg/ports/kubernetes"
	"github.com/notivenl/uptime-kubernetes/pkg/process"
	v1 "k8s.io/api/networking/v1"
)

func TestUptimeClient(t *testing.T) {
	tc := []struct {
		env         map[string]string
		expectPanic bool
	}{
		{
			env:         map[string]string{},
			expectPanic: true,
		},
		{
			env: map[string]string{
				"INPUT":  "notexisting",
				"OUTPUT": "uptimerobot",
			},
			expectPanic: true,
		},
		{
			env: map[string]string{
				"INPUT":  "annotation",
				"OUTPUT": "notexisting",
			},
			expectPanic: true,
		},
		{
			env: map[string]string{
				"INPUT":  "annotation",
				"OUTPUT": "uptimerobot",
			},
			expectPanic: false,
		},
	}

	for _, tt := range tc {
		// set env vars
		for k, v := range tt.env {
			os.Setenv(k, v)
		}

		func() {
			defer func() {
				if r := recover(); !tt.expectPanic && r != nil {
					t.Errorf("expected no panic: %s", r)
				}
			}()

			// create config
			conf := config.NewConfig()
			_ = client.NewUptimeClient(conf, &mockKubernetes.RestMock{
				GetInternalClientFunc: func() (kubernetes.Kubernetes, error) {
					return &mockKubernetes.KubernetesMock{
						GetIngressesFunc: func(ctx context.Context) ([]v1.Ingress, error) {
							return []v1.Ingress{
								{
									Spec: v1.IngressSpec{
										Rules: []v1.IngressRule{
											{
												Host: "test.com",
											},
										},
									},
								},
							}, nil
						},
					}, nil
				},
			})
		}()
	}
}

func TestUptimeClient_Loop(t *testing.T) {
	tc := []struct {
		input  *mock.InputMock
		output *mock.OutputMock
	}{
		{
			input: &mock.InputMock{
				InitFunc: func(configMoqParam *config.Config) error {
					return nil
				},
				ProcessFunc: func(ingress v1.Ingress, dataCh chan process.Data) {
					dataCh <- process.Data{
						Url: "test.com",
					}
				},
			},
			output: &mock.OutputMock{
				InitFunc: func(configMoqParam *config.Config) error {
					return nil
				},
				ProcessFunc: func(data process.Data) {
					if data.Url != "test.com" {
						t.Errorf("expected url to be test.com, got %s", data.Url)
					}
				},
			},
		},
	}

	for i, tt := range tc {
		// set env vars
		os.Setenv("INPUT", fmt.Sprintf("mock_input_%d", i))
		os.Setenv("OUTPUT", fmt.Sprintf("mock_output_%d", i))
		os.Setenv("TICKSPEED", "1")

		process.RegisterInput(fmt.Sprintf("mock_input_%d", i), tt.input)
		process.RegisterOutput(fmt.Sprintf("mock_output_%d", i), tt.output)

		conf := config.NewConfig()
		c := client.NewUptimeClient(conf, &mockKubernetes.RestMock{
			GetInternalClientFunc: func() (kubernetes.Kubernetes, error) {
				return &mockKubernetes.KubernetesMock{
					GetIngressesFunc: func(ctx context.Context) ([]v1.Ingress, error) {
						return []v1.Ingress{
							{
								Spec: v1.IngressSpec{
									Rules: []v1.IngressRule{
										{
											Host: "test.com",
										},
									},
								},
							},
						}, nil
					},
				}, nil
			},
		})

		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("expected no panic: %s", r)
				}
			}()
			go c.Start()
			time.Sleep(2 * time.Second)
			c.Stop()

			if len(tt.input.ProcessCalls()) == 0 {
				t.Errorf("expected input.Process to be called")
			}

			if len(tt.output.ProcessCalls()) == 0 {
				t.Errorf("expected input.Output to be called")
			}
		}()
	}
}
