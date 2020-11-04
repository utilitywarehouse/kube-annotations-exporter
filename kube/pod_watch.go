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

type podWatcher struct {
	annotations  []string
	client       kubernetes.Interface
	resyncPeriod time.Duration
	stopChannel  chan struct{}
	store        cache.Store
	Metrics      metrics.PrometheusInterface
}

func NewPodWatcher(client kubernetes.Interface, resyncPeriod time.Duration, metrics metrics.PrometheusInterface, annotations []string) *podWatcher {
	return &podWatcher{
		annotations:  annotations,
		client:       client,
		resyncPeriod: resyncPeriod,
		stopChannel:  make(chan struct{}),
		Metrics:      metrics,
	}
}

func (pw *podWatcher) updateMetrics() {
	podList := pw.List()
	pw.Metrics.UpdatePodAnnotations(podList, pw.annotations)
}

func (pw *podWatcher) eventHandler(eventType watch.EventType, old *v1.Pod, new *v1.Pod) {
	switch eventType {
	case watch.Added, watch.Modified, watch.Deleted:
		pw.updateMetrics()
	default:
		fmt.Println(
			fmt.Sprintf(
				"[Info] Unknown pod event received: %v",
				eventType,
			),
		)
	}
}

func (pw *podWatcher) Start() {
	listWatch := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return pw.client.CoreV1().Pods("").List(context.Background(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return pw.client.CoreV1().Pods("").Watch(context.Background(), options)
		},
	}
	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pw.eventHandler(watch.Added, nil, obj.(*v1.Pod))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pw.eventHandler(watch.Modified, oldObj.(*v1.Pod), newObj.(*v1.Pod))
		},
		DeleteFunc: func(obj interface{}) {
			pw.eventHandler(watch.Deleted, obj.(*v1.Pod), nil)
		},
	}
	store, controller := cache.NewInformer(listWatch, &v1.Pod{}, pw.resyncPeriod, eventHandler)
	pw.store = store
	fmt.Println("[Info] Starting pod watcher")
	// Running controller will block until writing on the stop channel.
	controller.Run(pw.stopChannel)
	fmt.Println("[Info] Stopped pod watcher")
}

func (pw *podWatcher) Stop() {
	fmt.Println("[Info] Stopping pod watcher")
	close(pw.stopChannel)
}

func (pw *podWatcher) List() []v1.Pod {
	var podList []v1.Pod
	for _, obj := range pw.store.List() {
		pod, ok := obj.(*v1.Pod)
		if !ok {
			fmt.Println(
				fmt.Sprintf(
					"[Error] Cannot read pod object: %s",
					obj,
				),
			)
			continue
		}
		podList = append(podList, *pod)
	}
	return podList
}
