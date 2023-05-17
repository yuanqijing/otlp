package receiver

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"k8s.io/klog"
)

const dataFormatProtobuf = "protobuf"

// Receiver is the type used to handle metrics from OpenTelemetry exporters.
type Receiver struct {
	//pmetricotlp.UnimplementedGRPCServer
	//nextConsumer consumer.Metrics
	//obsrecv      *obsreport.Receiver
}

func (r *Receiver) unexported() {
	//TODO implement me
	panic("implement me")
}

// New creates a new Receiver reference.
func New() *Receiver {
	return &Receiver{
		//nextConsumer: nextConsumer,
		//obsrecv:      obsrecv,
	}
}

// Export implements the service Export metrics func.
func (r *Receiver) Export(ctx context.Context, req pmetricotlp.ExportRequest) (pmetricotlp.ExportResponse, error) {
	md := req.Metrics()
	//dataPointCount := md.DataPointCount()
	//if dataPointCount == 0 {
	//	return pmetricotlp.NewExportResponse(), nil
	//}

	//ctx = r.obsrecv.StartMetricsOp(ctx)
	klog.Infof("req received: %s", spew.Sdump(md))
	//err := r.nextConsumer.ConsumeMetrics(ctx, md)
	//r.obsrecv.EndMetricsOp(ctx, dataFormatProtobuf, dataPointCount, err)

	return pmetricotlp.NewExportResponse(), nil
}
