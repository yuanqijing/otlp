package main

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/yuanqijing/otlp/testdata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
	"net"
)

func main() {
	md := testdata.GenerateMetrics(1)
	req := pmetricotlp.NewExportRequestFromMetrics(md)

	metricSink := new(consumertest.MetricsSink)
	metricsClient := makeMetricsServiceClient(metricSink)
	resp, err := metricsClient.Export(context.Background(), req)

	if err != nil {
		klog.Fatalf("Failed to export metrics: %v", err)
	}

	klog.Infof("resp: %s", spew.Sdump(resp))
}

func makeMetricsServiceClient(mc consumer.Metrics) pmetricotlp.GRPCClient {
	addr := otlpReceiverOnGRPCServer(mc)

	klog.Infof("addr: %s", addr.String())

	cc, err := grpc.Dial(addr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		klog.Fatalf("Failed to dial: %v", err)
	}

	return pmetricotlp.NewGRPCClient(cc)
}

func otlpReceiverOnGRPCServer(mc consumer.Metrics) net.Addr {
	ln, err := net.Listen("tcp", "localhost:")
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}

	set := receivertest.NewNopCreateSettings()
	set.ID = component.NewIDWithName("otlp", "metrics")
	obsrecv, err := obsreport.NewReceiver(obsreport.ReceiverSettings{
		ReceiverID:             set.ID,
		Transport:              "grpc",
		ReceiverCreateSettings: set,
	})
	r := New(mc, obsrecv)
	// Now run it as a gRPC server
	srv := grpc.NewServer()
	pmetricotlp.RegisterGRPCServer(srv, r)
	go func() {
		_ = srv.Serve(ln)
	}()

	return ln.Addr()
}
