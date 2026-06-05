package initialization

import (
	"context"
	"fmt"

	"go_template/internal/models"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ConnectS3(cfg models.Config) (*s3.Client, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.S3.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for S3: %w", err)
	}
	return s3.NewFromConfig(awsCfg), nil
}
