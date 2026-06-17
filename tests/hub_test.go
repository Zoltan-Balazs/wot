package tests

import (
	"testing"

	"github.com/Zoltan-Balazs/wot/src/hub"
)

func TestFindNearbyHubs_Budapest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Budapest Ferihegy airport is ~16 km SE of city centre
	hubs, err := hub.FindNearby(47.4979, 19.0402, 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(hubs) == 0 {
		t.Fatal("expected at least one hub near Budapest, got none")
	}
	for i := 1; i < len(hubs); i++ {
		if hubs[i].Distance < hubs[i-1].Distance {
			t.Fatalf("results not sorted: hubs[%d].Distance=%.2f < hubs[%d].Distance=%.2f",
				i, hubs[i].Distance, i-1, hubs[i-1].Distance)
		}
	}
}

func TestFindNearbyHubs_EmptyOcean(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Mid-Atlantic — no airports within 1 km
	hubs, err := hub.FindNearby(0, -30, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(hubs) != 0 {
		t.Fatalf("expected no hubs in mid-Atlantic, got %d", len(hubs))
	}
}
