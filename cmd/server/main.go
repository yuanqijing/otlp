package main

import (
	"github.com/yuanqijing/otlp/pkg/receiver"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"google.golang.org/grpc"
	"k8s.io/klog"
	"net"

	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	stopCtx := signals.SetupSignalHandler()

	var ln net.Listener
	var err error
	if ln, err = net.Listen("tcp", "localhost:16816"); err != nil {
		klog.Fatal(err)
	}

	klog.Infof("listening on addr: %s", ln.Addr().String())

	recv := receiver.New()

	srv := grpc.NewServer()
	pmetricotlp.RegisterGRPCServer(srv, recv)

	go func() {
		_ = srv.Serve(ln)
	}()

	<-stopCtx.Done()
}
