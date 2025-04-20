package config

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	AwsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func (cfg Config) LoadAwsConfig() aws.Config {
	conf, err := AwsConfig.LoadDefaultConfig(context.TODO(), AwsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.CF.ApiKey, cfg.CF.ApiSecret, "")), AwsConfig.WithRegion("auto"))

	if err != nil {
		log.Fatal().Msgf("Unable to load SDK config, %v", err)
	}

	log.Info().Msg("Success load AWS Config")
	return conf

}
