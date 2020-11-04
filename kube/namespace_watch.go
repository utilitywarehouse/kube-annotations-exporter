package kube

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/utilitywarehouse/kube-annotations-exporter/metrics"
)

type namespaceWatcher struct {
	annotations  []string
	client       kubernetes.Interface
	resyncPeriod time.Duration
	stopChannel  chan struct{}
	store        cache.Store
	Metrics      metrics.PrometheusInterface
}

func NewNamespaceWatcher(client kubernetes.Interface, resyncPeriod time.Duration, metrics metrics.PrometheusInterface, annotations []string) *namespaceWatcher {
	return &namespaceWatcher{
		annotations:  annotations,
		client:       client,
		resyncPeriod: resyncPeriod,
		stopChannel:  make(chan struct{}),
		Metrics:      metrics,
	}
}

func (nw *namespaceWatcher) updateNamespaceMetrics() {
	nsList := nw.List()
	nw.Metrics.UpdateNamespaceAnnotations(nsList, nw.annotations)
}

func (nw *namespaceWatcher) eventHandler(eventType watch.EventType, old *v1.Namespace, new *v1.Namespace) {
	switch eventType {
	case watch.Added, watch.Modified, watch.Deleted:
		nw.updateNamespaceMetrics()
	default:
		fmt.Println(
			fmt.Sprintf(
				"[Info] Unknown namespace event received: %v",
				eventType,
			),
		)
	}
}

func (nw *namespaceWatcher) Start() {
	listWatch := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return nw.client.CoreV1().Namespaces().List(context.Background(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return nw.client.CoreV1().Namespaces().Watch(context.Background(), options)
		},
	}
	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			nw.eventHandler(watch.Added, nil, obj.(*v1.Namespace))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			nw.eventHandler(watch.Modified, oldObj.(*v1.Namespace), newObj.(*v1.Namespace))
		},
		DeleteFunc: func(obj interface{}) {
			nw.eventHandler(watch.Deleted, obj.(*v1.Namespace), nil)
		},
	}
	store, controller := cache.NewInformer(listWatch, &v1.Namespace{}, nw.resyncPeriod, eventHandler)
	nw.store = store
	fmt.Println("[Info] Starting namespace watcher")
	// Running controller will block until writing on the stop channel.
	controller.Run(nw.stopChannel)
	fmt.Println("[Info] Stopped namespace watcher")
}

func (nw *namespaceWatcher) Stop() {
	fmt.Println("[Info] Stopping namespace watcher")
	close(nw.stopChannel)
}

func (nw *namespaceWatcher) List() []v1.Namespace {
	var nsList []v1.Namespace
	for _, obj := range nw.store.List() {
		ns, ok := obj.(*v1.Namespace)
		if !ok {
			fmt.Println(
				fmt.Sprintf(
					"[Error] Cannot read namespace object: %s",
					obj,
				),
			)
			continue
		}
		nsList = append(nsList, *ns)
	}
	return nsList
}
