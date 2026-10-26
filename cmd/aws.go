package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/viper"
)

func loadAWSConfig(ctx context.Context) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error
	if region := viper.GetString(keyRegion); region != "" {
		opts = append(opts, config.WithRegion(region))
	}
	if profile := viper.GetString(keyProfile); profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("loading AWS config: %w", err)
	}
	return cfg, nil
}
