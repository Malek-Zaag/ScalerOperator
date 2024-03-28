package controller

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func cluster_connect() (*rest.Config, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Println(kubeConfigPath)
	client, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		panic(err.Error())
	}
	return client, nil
}

func get_metrics(client *rest.Config) []v1beta1.PodMetrics {
	metrics, err := metrics.NewForConfig(client)
	if err != nil {
		panic(err.Error())
	}
	podMetrics, err := metrics.MetricsV1beta1().PodMetricses("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	// for _, item := range podMetrics.Items {
	// 	fmt.Printf("pod name is %s\n", item.GetName())
	// 	fmt.Printf("pod namesapce is %s\n", item.GetNamespace())
	// 	fmt.Printf("pod cpu usage is %d\n", item.Containers[0].Usage.Cpu().MilliValue())
	// 	fmt.Printf("pod memory usage is %d\n", item.Containers[0].Usage.Memory().Value()/(1024*1024))
	// }
	return podMetrics.Items
}
