package kafka

import "github.com/prometheus/client_golang/prometheus"

var (
	PublishTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "kafka",
		Subsystem: "producer",
		Name:      "publish_total",
		Help:      "Total events enqueued for publish.",
	})
	PublishDropped = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "kafka",
		Subsystem: "producer",
		Name:      "publish_dropped_total",
		Help:      "Events dropped because local queue was full.",
	})
	PublishErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "kafka",
		Subsystem: "producer",
		Name:      "publish_errors_total",
		Help:      "Errors while writing messages to Kafka.",
	})

	ConsumeTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "kafka",
		Subsystem: "consumer",
		Name:      "consume_total",
		Help:      "Total events consumed from Kafka.",
	})
	ConsumeErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "kafka",
		Subsystem: "consumer",
		Name:      "consume_errors_total",
		Help:      "Errors while consuming/decoding messages.",
	})
	ConsumerLag = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "kafka",
		Subsystem: "consumer",
		Name:      "lag",
		Help:      "Consumer lag (approximate).",
	})
)

func init() {
	prometheus.MustRegister(PublishTotal, PublishDropped, PublishErrors, ConsumeTotal, ConsumeErrors, ConsumerLag)
}
