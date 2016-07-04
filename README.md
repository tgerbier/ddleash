Datadog on a Leash
==================

Keep your Datadog on a leash by monitoring your metric usage.

ddleash is a set of Go packages and commands to help you interface
with Datadog, keep everything under control as well as debug how your
infrastructure uses the service.


Monitor your monitoring
-----------------------

ddleash includes the `ddleash-metrics-monitor` command for fetching
all metrics known by Datadog and record on datadog itself the total
number. You can then build a dashboard to keep an eye on this usage
and spot unusual spikes in metrics, alert your teams and avoid the
unpleasant bill.

This command takes a snapshot of the number of metrics existing in the
past day, and sends this value as a "gauge" metric to
Datadog. Typically, you would run the following command periodically
as a cron job.

Usage Example:

```
ddleash-metrics-monitor
```


Interface with Datadog programmatically
---------------------------------------

By importing these packages, you can write simple Go programs to
interface with Datadog for tasks unavailable with their public
API. For now, a few methods (used by `ddleash-metrics-monitor`) are
exposed allowing you to fetch summary information about a particular
metric.

For example, the following snippet will print all the metrics known by
Datadog in the last month:

```go
package main

import (
	"fmt"
	"os"

	"github.com/MattHauglustaine/ddleash"
)

func main() {
	ddleash, err := ddleash.New(ddleash.Config{
		Team:     "foo",
		Username: "bar",
		Password: "baz",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := ddleash.Login(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	metrics, err := ddleash.FetchAllMetricNames()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(len(metrics))
}
```
