package prometheuscustomexporter

import (
	"bytes"
	"context"
	"errors"

	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	// TODO: once this repository has been transferred to the
	// official census-ecosystem location, update this import path.

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/translator/internaldata"
)

var errBlankPrometheusAddress = errors.New("expecting a non-blank address to run the Prometheus metrics handler")

type prometheusCustomExporter struct {
	name         string
	exporter     *PrometheusContribExporter
	shutdownFunc func() error
}

func (pe *prometheusCustomExporter) Start(_ context.Context, _ component.Host) error {
	return nil
}

func (pe *prometheusCustomExporter) ConsumeMetrics(ctx context.Context, md pdata.Metrics) error {
	ocmds := internaldata.MetricsToOC(md)
	for _, ocmd := range ocmds {
		merged := make(map[string]*metricspb.Metric)
		for _, metric := range ocmd.Metrics {
			merge(merged, metric)
		}
		for _, metric := range merged {
			_ = pe.exporter.ExportMetric(ctx, ocmd.Node, ocmd.Resource, metric)
		}
	}
	return nil
}

// The underlying exporter overwrites timeseries when there are conflicting metric signatures.
// Therefore, we need to merge timeseries that share a metric signature into a single metric before sending.
func merge(m map[string]*metricspb.Metric, metric *metricspb.Metric) {
	key := metricSignature(metric)
	current, ok := m[key]
	if !ok {
		m[key] = metric
		return
	}
	current.Timeseries = append(current.Timeseries, metric.Timeseries...)
}

// Unique identifier of a given promtheus metric
// Assumes label keys are always in the same order
func metricSignature(metric *metricspb.Metric) string {
	var buf bytes.Buffer
	buf.WriteString(metric.GetMetricDescriptor().GetName())
	labelKeys := metric.GetMetricDescriptor().GetLabelKeys()
	for _, labelKey := range labelKeys {
		buf.WriteString("-" + labelKey.Key)
	}
	return buf.String()
}

// Shutdown stops the exporter and is invoked during shutdown.
func (pe *prometheusCustomExporter) Shutdown(context.Context) error {
	return pe.shutdownFunc()
}
