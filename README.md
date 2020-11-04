# kube-annotations-exporter

[![Docker Repository on Quay](https://quay.io/repository/utilitywarehouse/kube-annotations-exporter/status "Docker Repository on Quay")](https://quay.io/repository/utilitywarehouse/kube-annotations-exporter)

Exports pod and namespace annotations as metrics.

The namespace metrics are a substitute for old kube-state-metrics functionality:
- https://github.com/kubernetes/kube-state-metrics/pull/859

## Metrics

- `kube_namespace_annotations`: includes labels for `namespace`, `key` and `value`
- `kube_pod_annotations`: includes labels for `pod`, `namespace`, `key` and `value`

## Cardinality

By default, every annotation is exported for every pod and namespace. This
could produce a large number of series, so to mitigate this, you can provide a list
of annotations that you want to export with the flags `-namespace-annotations`
and `-pod-annotations`.

For example:
```
./kube-annotations-exporter \
  -pod-annotations="prometheus.io/scrape" \
  -pod-annotations="prometheus.io/path" \
  -pod-annotations="kubernetes.io/psp"
```

The flags can also be provided as a comma-delimited list:
```
./kube-annotations-exporter \
  -pod-annotations="prometheus.io/scrape,prometheus.io/path,kubernetes.io/psp"
```

