package main

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/yuanqijing/otlp/testdata"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"
)

func main() {
	cc, err := grpc.Dial("localhost:16816", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		klog.Fatalf("Failed to dial: %v", err)
	}

	client := pmetricotlp.NewGRPCClient(cc)
	md := testdata.GenerateMetrics(1)
	req := pmetricotlp.NewExportRequestFromMetrics(md)
	resp, err := client.Export(context.Background(), req)

	if err != nil {
		klog.Fatalf("Failed to export metrics: %v", err)
	}

	klog.Infof("resp: %s", spew.Sdump(resp))
}
