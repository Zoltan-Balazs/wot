package hub

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"

	"github.com/IBM/cloudant-go-sdk/cloudantv1"
	"github.com/Zoltan-Balazs/wot/src/cloudant"
	"github.com/Zoltan-Balazs/wot/src/config"
	"github.com/Zoltan-Balazs/wot/src/geo"
)

const searchLimit = 200

type Hub struct {
	Name     string
	Lat      float64
	Lon      float64
	Distance float64
}

func FindNearby(lat, lon, distKm float64) ([]Hub, error) {
	client, err := cloudant.Get()
	if err != nil {
		return nil, err
	}

	cfg := config.Get()
	minLat, maxLat, minLon, maxLon := geo.BoundingBox(lat, lon, distKm)
	query := fmt.Sprintf("lat:[%f TO %f] AND lon:[%f TO %f]", minLat, maxLat, minLon, maxLon)

	slog.Debug("querying Cloudant", "query", query)

	rows, err := fetchAllRows(client, cfg, query)
	if err != nil {
		return nil, err
	}

	var hubs []Hub
	for _, row := range rows {
		h, ok := rowToHub(row, lat, lon)
		if !ok {
			continue
		}
		if h.Distance <= distKm {
			hubs = append(hubs, h)
		}
	}

	slices.SortFunc(hubs, func(a, b Hub) int { return cmp.Compare(a.Distance, b.Distance) })
	slog.Info("found nearby hubs", "count", len(hubs), "lat", lat, "lon", lon, "distKm", distKm)
	return hubs, nil
}

func fetchAllRows(client *cloudantv1.CloudantV1, cfg config.Config, query string) ([]cloudantv1.SearchResultRow, error) {
	var allRows []cloudantv1.SearchResultRow
	var bookmark string

	for {
		opts := client.NewPostSearchOptions(cfg.DBName, cfg.DesignDoc, cfg.IndexName, query).
			SetLimit(searchLimit)
		if bookmark != "" {
			opts.SetBookmark(bookmark)
		}

		result, _, err := client.PostSearch(opts)
		if err != nil {
			return nil, fmt.Errorf("querying Cloudant: %w", err)
		}

		allRows = append(allRows, result.Rows...)

		if len(result.Rows) < searchLimit {
			break
		}
		if result.TotalRows != nil && len(allRows) >= int(*result.TotalRows) {
			break
		}
		if result.Bookmark == nil || *result.Bookmark == bookmark {
			break
		}

		bookmark = *result.Bookmark
	}

	return allRows, nil
}

func rowToHub(row cloudantv1.SearchResultRow, userLat, userLon float64) (Hub, bool) {
	lat, ok := row.Fields["lat"].(float64)
	if !ok {
		return Hub{}, false
	}
	lon, ok := row.Fields["lon"].(float64)
	if !ok {
		return Hub{}, false
	}
	name, ok := row.Fields["name"].(string)
	if !ok {
		return Hub{}, false
	}
	return Hub{
		Name:     name,
		Lat:      lat,
		Lon:      lon,
		Distance: geo.Haversine(userLat, userLon, lat, lon),
	}, true
}
