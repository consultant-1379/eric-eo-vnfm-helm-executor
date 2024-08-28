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
package main

import (
    "context"
    stderr "errors"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/client"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/config"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/errors"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/logging"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/mapper"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/message"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/model"
    "os"
    "plugin"
)

func main() {
    ctx := context.Background()

    conf, err := config.InitConfig()
    logging.InitLogger(conf)
    exitIfError(err, conf)

    redisClient := client.NewRedisClusterClient(ctx, conf)
    commandCtxString, err := redisClient.GetValue(conf.RedisKey)
    exitIfError(err, conf)
    logging.Log().Infof("Command contest: %s", commandCtxString)

    commandCtx := &model.CommandContext{}
    err = mapper.MapToStruct(commandCtxString, commandCtx)
    exitIfError(err, conf)

    decryptClient := client.NewDecryptClient(conf.CryptoHost)
    err = runHelmPlugin(commandCtx, decryptClient, redisClient)
    exitCode := errors.GetExitCode(err)

    camundaClient := client.NewCamundaClient(conf.WfsCamundaUrl)
    camundaMessage := message.BuildCamundaMessage(err, conf.ProcessInstanceId, commandCtx)
    logging.Log().Infof("Camunda message: %v", camundaMessage)
    camundaResponse, err := camundaClient.SendMessage(camundaMessage)
    if err != nil {
        logging.Log().Fatalf("An error occurred while sending camunda message: %s", err.Error())
    }
    logging.Log().Infof("Camunda response: %s", camundaResponse)
    os.Exit(exitCode)
}

func runHelmPlugin(commandCtx *model.CommandContext, decryptClient client.Decryptor, redisClient client.RedisClient) error {
    logging.Log().Infof("Helm version to be used: %s", commandCtx.HelmClientVersion)
    var helmPluginPath string
    switch commandCtx.HelmClientVersion {
    case "helm-3.10":
        helmPluginPath = "/usr/bin/helm-plugin-3.10.1.so"
    case "helm-3.12":
        helmPluginPath = "/usr/bin/helm-plugin-3.12.0.so"
    case "helm-3.13":
        helmPluginPath = "/usr/bin/helm-plugin-3.13.0.so"
    case "helm-3.14", "helm-latest":
        helmPluginPath = "/usr/bin/helm-plugin-3.14.2.so"
    default:
        helmPluginPath = "/usr/bin/helm-plugin-3.8.1.so"
    }

    p, err := plugin.Open(helmPluginPath)
    if err != nil {
        logging.Log().Errorf("An error occurred while opening helm plugin: %s", err.Error())
        return err
    }

    runFunc, err := p.Lookup("RunHelmCommand")
    if err != nil {
        logging.Log().Errorf("An error occurred while helm plugin lookup: %s", err.Error())
        return err
    }

    return runFunc.(func(*model.CommandContext, client.Decryptor, client.RedisClient) error)(commandCtx, decryptClient, redisClient)
}

func exitIfError(err error, conf *config.Config) {
    if err != nil {
        if conf.ProcessInstanceId == "" || conf.WfsCamundaUrl == "" {
            logging.Log().Fatalf("An error occurred: %s", err.Error())
        }
        logging.Log().Errorf("An error occcurred: %s", err.Error())
        camundaClient := client.NewCamundaClient(conf.WfsCamundaUrl)
        camundaMessage := message.BuildCamundaMessage(err, conf.ProcessInstanceId, nil)
        logging.Log().Infof("Camunda message: %v", camundaMessage)
        camundaResponse, responseErr := camundaClient.SendMessage(camundaMessage)
        if responseErr != nil {
            err = stderr.Join(err, responseErr)
            logging.Log().Fatalf("An error occurred while sending camunda message: %s", err.Error())
        }
        logging.Log().Infof("Camunda response: %s", camundaResponse)
        os.Exit(1)
    }
}
