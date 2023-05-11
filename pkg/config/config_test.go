package config_test

import (
	"os"
	"testing"

	"github.com/notivenl/uptime-kubernetes/pkg/config"
)

func Test_NewConfig(t *testing.T) {
	tc := []struct {
		env         map[string]string
		expected    config.Config
		expectPanic bool
	}{
		{
			env: map[string]string{
				"INPUT":     "test",
				"OUTPUT":    "test",
				"TICKSPEED": "1",
			},
			expected: config.Config{
				Input:     "test",
				Output:    "test",
				TickSpeed: 1,
			},
			expectPanic: false,
		},
		{
			env: map[string]string{
				"INPUT":     "",
				"OUTPUT":    "",
				"TICKSPEED": "false",
			},
			expected: config.Config{
				Input:     "",
				Output:    "",
				TickSpeed: 0,
			},
			expectPanic: true,
		},
		{
			env: map[string]string{
				"INPUT":     "",
				"OUTPUT":    "",
				"TICKSPEED": "",
			},
			expected: config.Config{
				Input:     "",
				Output:    "",
				TickSpeed: 0,
			},
			expectPanic: true,
		},
	}

	for _, tt := range tc {
		for k, v := range tt.env {
			os.Setenv(k, v)
		}

		func() {
			defer func() {
				if r := recover(); !tt.expectPanic && r != nil {
					t.Errorf("expected no panic: %s", r)
				}
			}()

			c := config.NewConfig()

			if c.Input != tt.expected.Input {
				t.Errorf("expected %s, got %s", tt.expected.Input, c.Input)
			}

			if c.Output != tt.expected.Output {
				t.Errorf("expected %s, got %s", tt.expected.Output, c.Output)
			}

			if c.TickSpeed != tt.expected.TickSpeed {
				t.Errorf("expected %d, got %d", tt.expected.TickSpeed, c.TickSpeed)
			}
		}()
	}
}

func Test_Config_Load(t *testing.T) {
	type testStruct struct {
		TestA string `env:"TEST"`
		TestB int
		TestC bool
		TestD struct { // structs are not supported
			TestE string
		}
		TestF *string // pointers are not supported
		TestG float32 // floats are not supported
	}

	env := map[string]string{
		"TEST":  "test",
		"TESTB": "1",
		"TESTC": "true",
		"TESTE": "test",
		"TESTF": "test",
		"TESTG": "1.1",
	}

	for k, v := range env {
		os.Setenv(k, v)
	}

	c := &config.Config{}

	conf := &testStruct{}

	c.Load(conf)

	if conf.TestA != "test" {
		t.Errorf("expected %s, got %s", "test", conf.TestA)
	}

	if conf.TestB != 1 {
		t.Errorf("expected %d, got %d", 1, conf.TestB)
	}

	if conf.TestC != true {
		t.Errorf("expected %t, got %t", true, conf.TestC)
	}

	if conf.TestD.TestE != "" {
		t.Errorf("expected %s, got %s", "", conf.TestD.TestE)
	}

	if conf.TestF != nil {
		t.Errorf("expected %s, got %s", "", *conf.TestF)
	}

	if conf.TestG != 0 {
		t.Errorf("expected %f, got %f", 1.1, conf.TestG)
	}
}
