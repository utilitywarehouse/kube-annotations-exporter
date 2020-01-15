package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/api/core/v1"
)

// PrometheusInterface allows for mocking out the functionality of Prometheus when testing.
type PrometheusInterface interface {
	UpdateNamespaceAnnotations([]v1.Namespace)
	DeleteNamespace(*v1.Namespace)
}

type Prometheus struct {
	namespaceAnnotations *prometheus.GaugeVec
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

	prometheus.MustRegister(p.namespaceAnnotations)
}

func (p *Prometheus) UpdateNamespaceAnnotations(nsList []v1.Namespace) {
	// Flush so annotations that no longer exist get deleted
	p.namespaceAnnotations.Reset()

	// Then set a metric for each of the existing annotations to 1
	for _, ns := range nsList {
		for key, value := range ns.Annotations {
			p.namespaceAnnotations.With(prometheus.Labels{
				"namespace": ns.Name,
				"key":       key,
				"value":     value,
			}).Set(1)
		}
	}

}

func (p *Prometheus) DeleteNamespace(ns *v1.Namespace) {
	// func (GaugeVec) Delete: returns true if a metric was deleted and
	// and false otherwise. No panic caused if metric doesn't exist
	p.namespaceAnnotations.Delete(prometheus.Labels{
		"namespace": ns.Name,
	})
}
