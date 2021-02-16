module github.com/yorikya/otel-custom-prometheus

go 1.14

require (
	github.com/census-instrumentation/opencensus-proto v0.3.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jaegerthrifthttpexporter v0.20.0 // indirect
	github.com/orijtech/prometheus-go-metrics-exporter v0.0.6
	github.com/prometheus/client_golang v1.9.0
	go.opentelemetry.io/collector v0.20.0
)
