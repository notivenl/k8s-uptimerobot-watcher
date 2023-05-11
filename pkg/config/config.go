package config

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// Config is the structure that defines and contains all the configuration
type Config struct {
	TickSpeed int
	Input     string
	Output    string
}

// NewConfig creates a new Config struct and loads the environment variables
func NewConfig() *Config {
	config := &Config{}

	config.Load(config)

	return config
}

// Load loads the environment variables into the Config struct automatically using reflection
func (c *Config) Load(config interface{}) {
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)

	t = t.Elem()
	v = v.Elem()

	// loop through all the fields in the struct
	for i := 0; i < t.NumField(); i++ {
		value := v.Field(i)

		// if the value is a pointer, get the value it points to
		// continue doing this until the value is not a pointer
		for value.Kind() == reflect.Pointer {
			value = value.Elem()
		}

		// only integers, strings or booleans are supported, continue if it is not one of those kinds
		if !slices.Contains([]reflect.Kind{
			reflect.Int,
			reflect.String,
			reflect.Bool,
		}, value.Kind()) {
			continue
		}

		// if an env tag is set, use that, otherwise use the field name in uppercase
		envName := t.Field(i).Tag.Get("env")
		if envName == "" {
			envName = strings.ToUpper(t.Field(i).Name)
		}

		if value.Kind() == reflect.Struct {
			// if the value is a struct, call Load on it
			c.Load(value.Addr().Interface())
		} else {
			// if the value is not a struct, set the value to the value from the environment
			loadedValue := getValueFromEnv(value.Type(), envName)
			if loadedValue.Kind() == reflect.Ptr {
				loadedValue = loadedValue.Elem()
			}

			if value.Kind() == reflect.Ptr {
				value.Set(loadedValue.Addr())
				continue
			}

			value.Set(loadedValue)
		}
	}
}

// getValueFromEnv gets the value from the environment and converts it to the correct type
// if the environment variable is not set, it returns the default value by creating a new value using its type
func getValueFromEnv(t reflect.Type, name string) reflect.Value {
	value := ""
	if v, ok := os.LookupEnv(name); ok {
		value = v
	} else {
		// return default value
		return reflect.New(t)
	}

	// convert value to correct type
	switch t.Kind() {
	case reflect.Int:
		v, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(v)
	case reflect.String:
		return reflect.ValueOf(value)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(v)
	default:
		// Only strings, ints and bools are supported
		// return default value
		return reflect.New(t)
	}
}
