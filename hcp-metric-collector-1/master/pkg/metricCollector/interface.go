package metricCollector

type Metric interface {
	GetMetric() []byte
}
