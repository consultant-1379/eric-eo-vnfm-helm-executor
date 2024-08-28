/*
 * COPYRIGHT Ericsson 2024
 *
 *
 *
 * The copyright to the computer program(s) herein is the property of
 *
 * Ericsson Inc. The programs may be used and/or copied only with written
 *
 * permission from Ericsson Inc. or in accordance with the terms and
 *
 * conditions stipulated in the agreement/contract under which the
 *
 * program(s) have been supplied.
 */
package client

import (
	"context"
	"gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/config"
	"gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/logging"
	"github.com/go-redis/redis/v8"
	"net"
	"time"
)

const retryAttempts = 3
const retryDelay = 5 * time.Second

type RedisClient interface {
	GetValue(string) (string, error)
}

type RedisClusterClient struct {
	ctx           context.Context
	rdb           *redis.ClusterClient
	retryAttempts int
	retryDelay    time.Duration
}

func NewRedisClusterClient(ctx context.Context, config *config.Config) *RedisClusterClient {
	return &RedisClusterClient{
		ctx:           ctx,
		rdb:           getClusterClient(config),
		retryAttempts: retryAttempts,
		retryDelay:    retryDelay,
	}
}

func (r *RedisClusterClient) GetValue(key string) (string, error) {
	var err error
	var stringValue string

	for i := 1; i <= r.retryAttempts; i++ {
		logging.Log().Infof("Retrieving value for key: %s, attempt %d of %d", key, i, r.retryAttempts)
		stringValue, err = r.rdb.Get(r.ctx, key).Result()

		if err == nil {
			return stringValue, nil
		}

		logging.Log().Errorf("An error occurred while retrieving value for key: %s. Error message: %s", key, err.Error())

		if i != r.retryAttempts {
			time.Sleep(r.retryDelay)
		}
	}

	return "", err
}

func getClusterClient(config *config.Config) *redis.ClusterClient {
	opts := &redis.ClusterOptions{
		Addrs:    []string{net.JoinHostPort(config.RedisHost, config.RedisPort)},
		Username: config.RedisUsername,
		Password: config.RedisPassword,
	}
	return redis.NewClusterClient(opts)
}
