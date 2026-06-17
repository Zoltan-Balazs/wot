package cli

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/Zoltan-Balazs/wot/src/hub"
)

type args struct {
	lat, lon, distKm float64
}

func parseArgs() (args, error) {
	if len(os.Args) != 4 {
		return args{}, fmt.Errorf("usage: %s <latitude> <longitude> <distance_km>\nexample: %s 47.5 19.0 50", os.Args[0], os.Args[0])
	}

	lat, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil || lat < -90 || lat > 90 {
		return args{}, fmt.Errorf("latitude must be a number between -90 and 90, got %q", os.Args[1])
	}

	lon, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil || lon < -180 || lon > 180 {
		return args{}, fmt.Errorf("longitude must be a number between -180 and 180, got %q", os.Args[2])
	}

	distKm, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil || distKm <= 0 {
		return args{}, fmt.Errorf("distance must be a positive number, got %q", os.Args[3])
	}

	return args{lat: lat, lon: lon, distKm: distKm}, nil
}

func Run() {
	a, err := parseArgs()
	if err != nil {
		slog.Error("invalid arguments", "error", err)
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	hubs, err := hub.FindNearby(a.lat, a.lon, a.distKm)
	if err != nil {
		slog.Error("query failed", "error", err)
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if len(hubs) == 0 {
		fmt.Printf("No transport hubs found within %.1f km of (%.4f, %.4f).\n", a.distKm, a.lat, a.lon)
		return
	}

	fmt.Printf("Transport hubs within %.1f km of (%.4f, %.4f):\n\n", a.distKm, a.lat, a.lon)

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "Name\tLatitude\tLongitude\tDistance")
	fmt.Fprintln(tw, "----\t--------\t---------\t--------")
	for _, h := range hubs {
		fmt.Fprintf(tw, "%s\t%.6f\t%.6f\t%.2f km\n", h.Name, h.Lat, h.Lon, h.Distance)
	}
	tw.Flush()
}
