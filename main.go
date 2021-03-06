package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/utilitywarehouse/kube-annotations-exporter/kube"
	"github.com/utilitywarehouse/kube-annotations-exporter/metrics"
)

var (
	flagNamespaceAnnotations = &StringSliceFlag{}
	flagPodAnnotations       = &StringSliceFlag{}
	flagKubeConfigPath       = flag.String("config", "", "Path of a kube config file, if not provided the app will try $KUBECONFIG, $HOME/.kube/config or in cluster config")
)

func main() {
	flag.Var(flagNamespaceAnnotations, "namespace-annotations", "Annotations to export for namespaces. Can be set multiple times and/or in comma-delimited form. By default all annotations will be exported.")
	flag.Var(flagPodAnnotations, "pod-annotations", "Annotations to export for pods. Can be set multiple times and/or in comma-delimited form. By default all annotations will be exported.")
	flag.Parse()

	metrics := &metrics.Prometheus{}
	metrics.Init()

	kubeClient, err := kube.GetClient(*flagKubeConfigPath)
	if err != nil {
		fmt.Printf("[Error] Cannot create kube client: %v", err)
		os.Exit(1)
	}

	nsWatcher := kube.NewNamespaceWatcher(
		kubeClient,
		metrics,
		flagNamespaceAnnotations.StringSlice(),
	)
	go nsWatcher.Start()

	podWatcher := kube.NewPodWatcher(
		kubeClient,
		metrics,
		flagPodAnnotations.StringSlice(),
	)
	go podWatcher.Start()

	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("[Error]: %v", http.ListenAndServe(":8080", nil))
}
