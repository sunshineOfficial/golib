package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func PrepareBucket(ctx context.Context, mc *minio.Client, bucket string) error {
	exists, err := mc.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("could not check if bucket exists: %w", err)
	}

	if exists {
		return nil
	}

	err = mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("could not create bucket %s: %w", bucket, err)
	}

	policy := Policy{
		Version: "2012-10-17",
		Statement: []Statement{
			{
				Effect: "Allow",
				Principal: Principal{
					AWS: []string{"*"},
				},
				Action:   []string{"s3:GetObject"},
				Resource: []string{fmt.Sprintf("arn:aws:s3:::%s/*", bucket)},
			},
		},
	}
	jsonPolicy, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("could not marshal policy: %w", err)
	}

	err = mc.SetBucketPolicy(ctx, bucket, string(jsonPolicy))
	if err != nil {
		return fmt.Errorf("could not set bucket policy: %w", err)
	}

	return nil
}

type Policy struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

type Statement struct {
	Effect    string    `json:"Effect"`
	Principal Principal `json:"Principal"`
	Action    []string  `json:"Action"`
	Resource  []string  `json:"Resource"`
}

type Principal struct {
	AWS []string `json:"AWS"`
}
