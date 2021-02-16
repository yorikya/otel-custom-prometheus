package prometheuscustomexporter

import (
	"context"

	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	"github.com/orijtech/prometheus-go-metrics-exporter"

	commonpb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/common/v1"
	resourcepb "github.com/census-instrumentation/opencensus-proto/gen-go/resource/v1"
)

type PrometheusContribExporter struct {
	*prometheus.Exporter
}

// ExportMetric is the method that the exporter uses to convert OpenCensus Proto-Metrics to Prometheus metrics.
func (exp *PrometheusContribExporter) ExportMetric(ctx context.Context, node *commonpb.Node, rsc *resourcepb.Resource, metric *metricspb.Metric) error {
	if metric == nil || len(metric.Timeseries) == 0 {
		return nil
	}

	for k, v := range rsc.Labels {
		metric.MetricDescriptor.LabelKeys = append(metric.MetricDescriptor.LabelKeys, &metricspb.LabelKey{
			Key: k,
		})

		metric.Timeseries[0].LabelValues = append(metric.Timeseries[0].LabelValues, &metricspb.LabelValue{
			Value: v,
		})
	}

	return exp.Exporter.ExportMetric(ctx, node, rsc, metric)
}

func NewPrometheusExporter(opts prometheus.Options) (*PrometheusContribExporter, error) {
	p, _ := prometheus.New(opts)
	return &PrometheusContribExporter{
		Exporter: p,
	}, nil
}
