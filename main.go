package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/utilitywarehouse/kube-namespace-annotations-exporter/kube"
	"github.com/utilitywarehouse/kube-namespace-annotations-exporter/metrics"
)

var (
	flagKubeConfigPath = flag.String("config", "", "Path of a kube config file, if not provided the app will try to get in cluster config")
	flagResyncPeriod   = flag.Int("resync-period", 60, "Namespace watcher cache resync period in minutes. Defaults to 60")
)

func main() {
	metrics := &metrics.Prometheus{}
	metrics.Init()

	flag.Parse()
	kubeClient, err := kube.GetClient(*flagKubeConfigPath)
	if err != nil {
		fmt.Printf("[Error] Cannot create kube client: %v", err)
		os.Exit(1)
	}

	resyncPeriod := time.Duration(*flagResyncPeriod) * time.Minute

	nsWatcher := kube.NewNamespaceWatcher(
		kubeClient,
		// Resync will trigger an onUpdate event for everything that is
		// stored in cache.
		resyncPeriod,
		metrics,
	)
	go nsWatcher.Start()

	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("[Error]: %v", http.ListenAndServe(":8080", nil))
}
