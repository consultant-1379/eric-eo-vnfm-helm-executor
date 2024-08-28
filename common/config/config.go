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
package config

import (
    "github.com/go-playground/validator/v10"
    "os"
)

type Config struct {
    RedisHost         string `validate:"required"`
    RedisPort         string `validate:"required"`
    RedisUsername     string `validate:"required"`
    RedisPassword     string `validate:"required"`
    RedisKey          string `validate:"required"`
    CryptoHost        string `validate:"required"`
    WfsCamundaUrl     string `validate:"required"`
    ProcessInstanceId string `validate:"required"`
    LogstashHost      string
    LogstashPort      string
    HelmDebug         string
    TraceId           string
}

func InitConfig() (*Config, error) {

    conf := &Config{
        RedisHost:         os.Getenv("REDIS_HOST"),
        RedisPort:         os.Getenv("REDIS_PORT"),
        RedisUsername:     os.Getenv("REDIS_USERNAME"),
        RedisPassword:     os.Getenv("REDIS_PASSWORD"),
        RedisKey:          os.Getenv("REDIS_KEY"),
        CryptoHost:        os.Getenv("CRYPTO_HOST"),
        WfsCamundaUrl:     os.Getenv("WFS_CAMUNDA_URL"),
        ProcessInstanceId: os.Getenv("PROCESS_INSTANCE_ID"),
        LogstashHost:      os.Getenv("LOGSTASH_HOST"),
        LogstashPort:      os.Getenv("LOGSTASH_PORT"),
        HelmDebug:         os.Getenv("HELM_DEBUG"),
        TraceId:           os.Getenv("TRACE_ID"),
    }

    validate := validator.New()
    if err := validate.Struct(conf); err != nil {
        return conf, err
    }

    return conf, nil
}
