package initialization

import (
	"context"
	"fmt"

	"go_template/internal/models"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

func ConnectSES(cfg models.Config) (*sesv2.Client, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.SES.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	return sesv2.NewFromConfig(awsCfg), nil
}
