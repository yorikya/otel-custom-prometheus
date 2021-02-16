package prometheuscustomexporter

import (
	"github.com/prometheus/client_golang/prometheus"

	"go.opentelemetry.io/collector/config/configmodels"
)

// Config defines configuration for Prometheus exporter.
type Config struct {
	configmodels.ExporterSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.

	// The address on which the Prometheus scrape handler will be run on.
	Endpoint string `mapstructure:"endpoint"`

	// Namespace if set, exports metrics under the provided value.
	Namespace string `mapstructure:"namespace"`

	// ConstLabels are values that are applied for every exported metric.
	ConstLabels prometheus.Labels `mapstructure:"const_labels"`

	// SendTimestamps will send the underlying scrape timestamp with the export
	SendTimestamps bool `mapstructure:"send_timestamps"`
}
