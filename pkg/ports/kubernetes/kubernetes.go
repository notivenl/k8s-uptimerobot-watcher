package kubernetes

import (
	"context"

	v1i "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//go:generate moq -pkg mock -out ../../../mocks/kubernetes/kubernetes.go . Kubernetes
//go:generate moq -pkg mock -out ../../../mocks/kubernetes/rest.go . Rest

// Kubernetes interface is used to mock the Kubernetes client for testing
type Kubernetes interface {
	GetIngresses(ctx context.Context) ([]v1i.Ingress, error)
}

// Rest interface is used to mock the Kubernetes client for testing
type Rest interface {
	GetInternalClient() (Kubernetes, error)
}

// default

// RestClient is the default implementation of the Kubernetes client
type RestClient struct {
}

// KubernetesClient is the default implementation of the Kubernetes client
type KubernetesClient struct {
	Client *kubernetes.Clientset
}

// NewRestClient returns a new instance of the default Kubernetes client
func NewRestClient() Rest {
	return &RestClient{}
}

// GetInternalClient gets the InClusterConfig from the Kubernetes API and returns a new Kubernetes client using it
func (r *RestClient) GetInternalClient() (Kubernetes, error) {
	conf, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, err
	}

	return &KubernetesClient{
		Client: client,
	}, nil
}

// GetIngresses returns all ingresses in all namespaces
func (k *KubernetesClient) GetIngresses(ctx context.Context) ([]v1i.Ingress, error) {
	namespaces, err := k.Client.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	ingresses := make([]v1i.Ingress, 0)

	// loop over all namespaces in search for ingresses
	for _, namespace := range namespaces.Items {
		ingressList, err := k.Client.NetworkingV1().Ingresses(namespace.Name).List(ctx, v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		ingresses = append(ingresses, ingressList.Items...)
	}

	return ingresses, nil
}
