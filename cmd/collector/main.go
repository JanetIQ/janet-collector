package main

import (
	"log"

	janetk8sreceiver "github.com/janetIQ/k8sreceiver"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/debugexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
)

func main() {
	info := component.BuildInfo{
		Command:     "janet-collector",
		Description: "JanetIQ custom OTel collector",
		Version:     "0.1.0",
	}

	set := otelcol.CollectorSettings{
		BuildInfo: info,
		Factories: components,
	}

	cmd := otelcol.NewCommand(set)
	if err := cmd.Execute(); err != nil {
		log.Fatalf("collector error: %v", err)
	}
}

// components registers all receivers, processors, and exporters
// that this collector binary supports.
func components() (otelcol.Factories, error) {
	var err error
	factories := otelcol.Factories{}

	// Receivers
	factories.Receivers, err = otelcol.MakeFactoryMap(
		janetk8sreceiver.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	// Processors
	factories.Processors, err = otelcol.MakeFactoryMap(
		batchprocessor.NewFactory(),
		memorylimiterprocessor.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	// Exporters
	factories.Exporters, err = otelcol.MakeFactoryMap(
		otlpexporter.NewFactory(),
		debugexporter.NewFactory(), // useful during local dev
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	return factories, nil
}
