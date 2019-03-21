package common

import (
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

func InitTracing(gcpProjectID string) error {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{ProjectID: gcpProjectID})
	if err != nil {
		return err
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	return nil
}
