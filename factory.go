package prometheuscustomexporter

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/orijtech/prometheus-go-metrics-exporter"
	"go.opentelemetry.io/collector/component"

	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr = "prometheuscustom"
)

// NewFactory creates a factory for OTLP exporter.
func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithMetrics(createMetricsExporter))
}

func createDefaultConfig() configmodels.Exporter {
	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		ConstLabels:    map[string]string{},
		SendTimestamps: false,
	}
}

func createMetricsExporter(
	_ context.Context,
	_ component.ExporterCreateParams,
	cfg configmodels.Exporter,
) (component.MetricsExporter, error) {
	pcfg := cfg.(*Config)

	addr := strings.TrimSpace(pcfg.Endpoint)
	if addr == "" {
		return nil, errBlankPrometheusAddress
	}

	opts := prometheus.Options{
		Namespace:      pcfg.Namespace,
		ConstLabels:    pcfg.ConstLabels,
		SendTimestamps: pcfg.SendTimestamps,
	}
	pe, err := NewPrometheusExporter(opts)
	if err != nil {
		return nil, err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// The Prometheus metrics exporter has to run on the provided address
	// as a server that'll be scraped by Prometheus.
	mux := http.NewServeMux()
	mux.Handle("/metrics", pe)

	srv := &http.Server{Handler: mux}
	go func() {
		_ = srv.Serve(ln)
	}()

	pexp := &prometheusCustomExporter{
		name:         cfg.Name(),
		exporter:     pe,
		shutdownFunc: ln.Close,
	}

	return pexp, nil
}
