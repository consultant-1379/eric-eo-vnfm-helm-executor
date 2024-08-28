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
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/errors"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/logging"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/model"
    "helm.sh/helm/v3/pkg/action"
    "helm.sh/helm/v3/pkg/cli"
    "time"
)

func (ex *executor) helmInstall(installParams *model.InstallParams) error {
    helmSettings, err := ex.getHelmSettings(&installParams.BaseParams)
    if err != nil {
        return err
    }

    installAction, err := ex.configureInstallAction(helmSettings, installParams)
    if err != nil {
        return err
    }

    chartName := getChartName(installParams)
    chartPath, err := installAction.ChartPathOptions.LocateChart(chartName, helmSettings)
    if err != nil {
        logging.Log().Errorf("An error occurred while downloading helm chart: %s. Error message: %s", chartName, err.Error())
        return errors.NewHelmError(err.Error())
    }

    chart, err := getChart(chartPath)
    if err != nil {
        return err
    }

    mergedValues, err := ex.mergeValues(helmSettings, installParams)
    if err != nil {
        return err
    }

    logging.Log().Infof("Starting helm install command for release %s in namespace %s",
        installParams.ReleaseName, installParams.Namespace)
    _, err = installAction.Run(chart, mergedValues)
    if err != nil {
        logging.Log().Errorf("An error occurred while running helm install command. Error message: %s", err.Error())
        return errors.NewHelmError(err.Error())
    }
    return nil

}

func (ex *executor) configureInstallAction(helmSettings *cli.EnvSettings, installParams *model.InstallParams) (*action.Install, error) {
    actionConfig, err := ex.getActionConfig(helmSettings)
    if err != nil {
        return nil, err
    }

    installAction := action.NewInstall(actionConfig)
    installAction.Namespace = installParams.Namespace
    installAction.ReleaseName = installParams.ReleaseName
    installAction.Timeout = time.Second * time.Duration(installParams.Timeout)
    installAction.RepoURL = installParams.ChartRepo
    installAction.DisableOpenAPIValidation = installParams.DisableOpenAPIValidation
    installAction.Wait = installParams.Wait
    installAction.DisableHooks = installParams.DisableHooks
    installAction.Atomic = installParams.Atomic
    installAction.CreateNamespace = installParams.CreateNamespace

    return installAction, nil
}
