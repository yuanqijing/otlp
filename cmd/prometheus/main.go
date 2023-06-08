package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"k8s.io/klog"
	"log"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"time"
)

type Listener struct {
	exporter      *prometheus.Exporter
	meterProvider *sdkmetric.MeterProvider
}

func main() {
	// Wait a bit to allow metrics to be collected
	stopCtx := signals.SetupSignalHandler()

	var l Listener
	var err error

	// Create Prometheus Exporter
	l.exporter, err = prometheus.New()
	if err != nil {
		log.Fatalf("Failed to initialize Prometheus exporter: %v", err)
	}

	// Create MeterProvider with Prometheus Exporter as Reader
	l.meterProvider = sdkmetric.NewMeterProvider(sdkmetric.WithReader(l.exporter))

	// Set MeterProvider in global OTel API

	// Use the meterProvider to get a named meter
	meter := l.meterProvider.Meter("my.instrumentation.library")

	// Create an integer counter instrument
	counter, err := meter.Int64Counter("my.counter")

	// Record a measurement
	ctx := context.Background()
	go func() {
		counters := 1
		for {
			select {
			case <-stopCtx.Done():
				return
			case <-time.After(time.Second * 1):
			}
			counter.Add(ctx, 1, otelmetric.WithAttributes(attribute.String("key", fmt.Sprintf("value-%d", counters))))
			counters++
			if counters > 10 {
				counter, _ = meter.Int64Counter("my.counter")
			}
		}
	}()
	counter.Add(ctx, 1, otelmetric.WithAttributes(attribute.String("key", fmt.Sprintf("value-%d", 3))))

	// Provide your own logic to get the listening address and the metrics path.
	listenAddress := ":8080"
	metricsPath := "/metrics"

	// Start serving metrics endpoint
	go func() {
		mux := http.NewServeMux()
		mux.Handle(metricsPath, promhttp.Handler())
		klog.Infof("listening on addr: %s", listenAddress)
		err := http.ListenAndServe(listenAddress, mux)
		if err != nil {
			log.Fatalf("Error serving HTTP: %v", err)
			return
		}
	}()

	<-stopCtx.Done()
	log.Println("Finished recording metrics")
}
