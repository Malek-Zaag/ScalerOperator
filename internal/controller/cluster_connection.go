package controller

import (
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func cluster_connect() (*rest.Config, error) {
	client, err := clientcmd.BuildConfigFromFlags("", "/mnt/c/Users/mk/.kube/config")

	if err != nil {
		panic(err.Error())
	}
	return client, nil
}

func get_metrics(client *rest.Config) {
	metrics, err := metrics.NewForConfig(client)
	if err != nil {
		panic(err.Error())
	}
	fmt.Print(metrics)
}
