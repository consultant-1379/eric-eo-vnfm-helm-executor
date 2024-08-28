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
package helm

import (
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/client"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/errors"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/fsutil"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/logging"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/model"
    "helm.sh/helm/v3/pkg/action"
    "helm.sh/helm/v3/pkg/cli"
    "helm.sh/helm/v3/pkg/cli/values"
    "helm.sh/helm/v3/pkg/getter"
    "os"
)

type executor struct {
    decryptClient client.Decryptor
    redisClient   client.RedisClient
}

func newExecutor(decryptClient client.Decryptor, redisClient client.RedisClient) *executor {
    return &executor{
        decryptClient: decryptClient,
        redisClient:   redisClient,
    }
}

func (ex *executor) getHelmSettings(baseParams *model.BaseParams) (*cli.EnvSettings, error) {
    logging.Log().Info("Configuring helm settings")
    configFilePath, err := ex.getFilePath(baseParams.ClusterConfigFileKey, "config")
    if err != nil {
        return nil, err
    }
    helmSettings := cli.New()
    helmSettings.KubeConfig = configFilePath
    helmSettings.SetNamespace(baseParams.Namespace)
    helmSettings.Debug = baseParams.HelmDebug
    return helmSettings, nil
}

func (ex *executor) getActionConfig(helmSettings *cli.EnvSettings) (*action.Configuration, error) {
    logging.Log().Info("Configuring helm action")
    debugFunc := logging.Log().Debugf
    actionConfig := new(action.Configuration)

    if err := actionConfig.Init(helmSettings.RESTClientGetter(), helmSettings.Namespace(), os.Getenv("HELM_DRIVER"), debugFunc); err != nil {
        logging.Log().Errorf("An error occurred while configuring helm action. Error message: %s", err.Error())
        return nil, errors.NewHelmError(err.Error())
    }
    return actionConfig, nil
}

func (ex *executor) getFilePath(redisKey, pattern string) (string, error) {
    encryptedContent, err := ex.redisClient.GetValue(redisKey)
    if err != nil {
        return "", err
    }

    decryptedContent, err := ex.decryptClient.Decrypt(&client.DecryptRequest{Ciphertext: encryptedContent})
    if err != nil {
        return "", err
    }

    tempFilePath, err := fsutil.SaveTempFile(decryptedContent, pattern)
    if err != nil {
        return "", err
    }
    return tempFilePath, nil
}

func (ex *executor) mergeValues(helmSettings *cli.EnvSettings, installParams *model.InstallParams) (map[string]interface{}, error) {
    opts := &values.Options{}
    if installParams.DayZeroValuesFileKey != "" {
        dayZeroValuesFilePath, err := ex.getFilePath(installParams.DayZeroValuesFileKey, "dayZeroValues")
        if err != nil {
            return nil, err
        }
        opts.ValueFiles = append(opts.ValueFiles, dayZeroValuesFilePath)
    }

    if installParams.ValuesFileKey != "" {
        valuesFilePath, err := ex.getFilePath(installParams.ValuesFileKey, "values")
        if err != nil {
            return nil, err
        }
        opts.ValueFiles = append(opts.ValueFiles, valuesFilePath)
    }

    if installParams.SetFlagValues != nil {
        opts.Values = installParams.SetFlagValues
    }

    provider := getter.All(helmSettings)
    logging.Log().Info("Merging user provided values")
    mergedValues, err := opts.MergeValues(provider)
    if err != nil {
        logging.Log().Errorf("An error occurred while merging values. Error message: %s", err.Error())
        return nil, errors.NewHelmError(err.Error())
    }

    return mergedValues, nil
}
