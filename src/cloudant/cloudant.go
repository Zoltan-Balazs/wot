package cloudant

import (
	"fmt"
	"log/slog"

	"github.com/IBM/cloudant-go-sdk/cloudantv1"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/Zoltan-Balazs/wot/src/config"
)

var client *cloudantv1.CloudantV1

func Get() (*cloudantv1.CloudantV1, error) {
	if client != nil {
		return client, nil
	}

	cfg := config.Get()
	svc, err := cloudantv1.NewCloudantV1(&cloudantv1.CloudantV1Options{
		URL:           cfg.CloudantURL,
		Authenticator: &core.NoAuthAuthenticator{},
	})
	if err != nil {
		return nil, fmt.Errorf("creating Cloudant client: %w", err)
	}

	slog.Info("cloudant client initialised", "url", cfg.CloudantURL)
	client = svc
	return client, nil
}
