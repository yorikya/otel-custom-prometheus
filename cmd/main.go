package main

import (
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/service"
	"go.opentelemetry.io/collector/service/defaultcomponents"

	jager "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jaegerthrifthttpexporter"
	prometheuscustom "github.com/yorikya/otel-custom-prometheus"
)

func main() {
	factories, err := components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
	}

	info := component.ApplicationStartInfo{
		ExeName:  "otelcol-filternaninf",
		LongName: "OpenTelemetry Collector with prometheus extend exporter",
		Version:  "1.0.0",
	}

	app, err := service.New(service.Parameters{ApplicationStartInfo: info, Factories: factories})
	if err != nil {
		log.Fatal("failed to construct the application: %w", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal("application run finished with error: %w", err)
	}
}

func components() (component.Factories, error) {
	var errs []error
	factories, err := defaultcomponents.Components()
	if err != nil {
		return component.Factories{}, err
	}

	exporters := []component.ExporterFactory{
		prometheuscustom.NewFactory(),
		jager.NewFactory(),
	}
	for _, ex := range factories.Exporters {
		exporters = append(exporters, ex)
	}
	factories.Exporters, err = component.MakeExporterFactoryMap(exporters...)
	if err != nil {
		errs = append(errs, err)
	}

	return factories, componenterror.CombineErrors(errs)
}
