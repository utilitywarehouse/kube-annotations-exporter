package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/api/core/v1"
)

// PrometheusInterface allows for mocking out the functionality of Prometheus when testing.
type PrometheusInterface interface {
	UpdateNamespaceAnnotations([]v1.Namespace, []string)
	UpdatePodAnnotations([]v1.Pod, []string)
}

type Prometheus struct {
	namespaceAnnotations *prometheus.GaugeVec
	podAnnotations       *prometheus.GaugeVec
}

// Init creates and registers the metrics.
func (p *Prometheus) Init() {

	p.namespaceAnnotations = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kube_namespace_annotations",
		Help: "Kubernetes namespace annotations",
	},
		[]string{
			// Namespace in question
			"namespace",
			// Annotation key
			"key",
			// Annotation value
			"value",
		},
	)

	p.podAnnotations = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kube_pod_annotations",
		Help: "Kubernetes pod annotations",
	},
		[]string{
			// Pod in question
			"pod",
			// Pod namespace
			"namespace",
			// Annotation key
			"key",
			// Annotation value
			"value",
		},
	)

	prometheus.MustRegister(p.namespaceAnnotations)
	prometheus.MustRegister(p.podAnnotations)
}

func (p *Prometheus) UpdateNamespaceAnnotations(nsList []v1.Namespace, annotations []string) {
	// Flush so annotations that no longer exist get deleted
	p.namespaceAnnotations.Reset()

	// Then set a metric for each of the existing annotations to 1
	for _, ns := range nsList {
		for key, value := range ns.Annotations {
			if len(annotations) == 0 || contains(annotations, key) {
				p.namespaceAnnotations.With(prometheus.Labels{
					"namespace": ns.Name,
					"key":       key,
					"value":     value,
				}).Set(1)
			}
		}
	}

}

func (p *Prometheus) UpdatePodAnnotations(podList []v1.Pod, annotations []string) {
	// Flush so annotations that no longer exist get deleted
	p.podAnnotations.Reset()

	// Then set a metric for each of the existing annotations to 1
	for _, pod := range podList {
		for key, value := range pod.Annotations {
			if len(annotations) == 0 || contains(annotations, key) {
				p.podAnnotations.With(prometheus.Labels{
					"pod":       pod.Name,
					"namespace": pod.Namespace,
					"key":       key,
					"value":     value,
				}).Set(1)
			}
		}
	}

}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}
