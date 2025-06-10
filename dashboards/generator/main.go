package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
)

var (
	dir = flag.String("dir", "", "output dir")
)

func main() {
	flag.Parse()

	for _, b := range []*dashboard.DashboardBuilder{
		NebiusDiskUserStats,
		NebiusGPU,
		NebiusObjectStorage,
		NebiusSharedFilesystem,
		NebiusLogging,
	} {
		d, err := b.Build()
		if err != nil {
			panic(err)
		}

		data, err := json.MarshalIndent(d, "", "  ")
		if err != nil {
			panic(err)
		}

		path := filepath.Join(*dir,
			fmt.Sprintf("%s.json", *d.Uid),
		)

		err = os.WriteFile(path, data, 0644)
		if err != nil {
			panic(err)
		}
	}
}
