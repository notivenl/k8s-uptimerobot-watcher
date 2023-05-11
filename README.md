
# Uptime Kubernetes

This program is a simple monitor creator for Kubernetes. It checks through ingresses in order to find what it should monitor and sends the data to a monitor service such as Uptimerobot to create a new monitor.

## How to use

This program can be installed either by creating the Kubernetes manifests manually, or by using the Helm chart.
In order to do it manually, see the [example folder](example).

For the Helm chart, see the helm charts over [here](https://github.com/notivenl/k8s-uptimerobot-chart).

## Configuration

The configuration is done through environment variables. The following variables are available:

| Variable | Description | Default |
| --- | --- | --- |
| `INPUT` | The input to use. This gets the data from all the registered ingresses, currently only `annotation` is availabe. | `nil` |
| `OUTPUT` | The output to use. This sends the data to a monitor service, currently only `uptimerobot` is available. | `nil` |
| `TICKSPEED` | The speed at which the program should check for changes in seconds. | `5` |

### Input/Output specific configuration
| `UPTIMEROBOT_API_TOKEN` | The API key to use for Uptimerobot. | `nil` |
