package main

import (
	"context"
	"fmt"
	"github.com/prometheus/prometheus/model/histogram"
	"math"
	"os"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	. "github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
)

func Example() {
	// Create a random dir to work in.  Open() doesn't require a pre-existing dir, but
	// we want to make sure not to make a mess where we shouldn't.
	dir, err := os.MkdirTemp("", "tsdb-test")
	noErr(err)

	// Open a TSDB for reading and/or writing.
	db, err := Open(dir, nil, nil, DefaultOptions(), nil)
	noErr(err)

	// Open an appender for writing.
	app := db.Appender(context.Background())

	series := labels.FromStrings("foo", "bar")

	// Ref is 0 for the first append since we don't know the reference for the series.
	_, err = app.Append(0, series, time.Now().Unix(), 123)
	noErr(err)

	// Another append for a second later.
	// Re-using the ref from above since it's the same series, makes append faster.
	time.Sleep(time.Second)

	series = labels.FromStrings("foo22", "barsfsfsdfdsfsfs", "foo", "bar")
	_, err = app.Append(0, series, time.Now().Unix(), 124)

	// Commit to storage.

	err = app.Commit()
	noErr(err)

	// In case you want to do more appends after app.Commit(),
	// you need a new appender.
	app = db.Appender(context.Background())
	his := &histogram.Histogram{
		Count:         10 + uint64(1*8),
		ZeroCount:     2 + uint64(1),
		ZeroThreshold: 0.001,
		Sum:           18.4 * float64(1+1),
		Schema:        1,
		PositiveSpans: []histogram.Span{
			{Offset: 0, Length: 2},
			{Offset: 1, Length: 2},
		},
		PositiveBuckets: []int64{int64(1 + 1), 1, -1, 0},
		NegativeSpans: []histogram.Span{
			{Offset: 0, Length: 2},
			{Offset: 1, Length: 2},
		},
		NegativeBuckets: []int64{int64(1 + 1), 1, -1, 0},
	}
	series = labels.FromStrings("foo22", "jjbjbjbjbjbjbjj", "foo", "bar")
	_, err = app.AppendHistogram(0, series, time.Now().Unix(), his, nil)
	err = app.Commit()

	//series = labels.FromStrings("foo22", "barsfsfsdfdsfsfs", "foo", "bar")
	//_, err = app.Append(0, series, time.Now().Unix(), 124)
	//noErr(err)
	//
	//err = app.Commit()
	//noErr(err)
	// ... adding more samples.

	querier, err := db.Querier(context.Background(), math.MinInt64, math.MaxInt64)
	noErr(err)
	ss := querier.Select(false, nil, labels.MustNewMatcher(labels.MatchEqual, "foo", "bar"))

	for ss.Next() {
		series := ss.At()
		fmt.Println("series:", series.Labels().String())

		it := series.Iterator(nil)
		for it.Next() == chunkenc.ValFloat {
			_, v := it.At() // We ignore the timestamp here, only to have a predictable output we can test against (below)
			fmt.Println("sample", v)
		}

		fmt.Println("it.Err():", it.Err())
	}
	fmt.Println("ss.Err():", ss.Err())
	ws := ss.Warnings()
	if len(ws) > 0 {
		fmt.Println("warnings:", ws)
	}
	err = querier.Close()
	noErr(err)

	// Clean up any last resources when done.
	err = db.Close()
	noErr(err)
	err = os.RemoveAll(dir)
	noErr(err)

	// Output:
	// series: {foo="bar"}
	// sample 123
	// sample 124
	// it.Err(): <nil>
	// ss.Err(): <nil>
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	Example()
}
