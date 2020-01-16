# kube-namespace-annotations-exporter

Exports namespace metrics as annotations.

Substitute for kube-state-metrics functionality before:
- https://github.com/kubernetes/kube-state-metrics/pull/859

# Metrics

- `kube_namespace_annotations`: includes labels for `namespace`, `key` and `value`
