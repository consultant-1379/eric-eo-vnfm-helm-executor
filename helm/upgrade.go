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
    "helm.sh/helm/v3/pkg/storage/driver"
    "time"
)

func (ex *executor) helmUpgrade(upgradeParams *model.UpgradeParams) error {
    helmSettings, err := ex.getHelmSettings(&upgradeParams.BaseParams)
    if err != nil {
        return err
    }

    upgradeAction, err := ex.configureUpgradeAction(helmSettings, upgradeParams)
    if err != nil {
        return err
    }
    if upgradeAction.Install {
        historyAction, err := ex.configureHistoryAction(helmSettings, &upgradeParams.BaseParams)
        if err != nil {
            return err
        }
        if _, err := historyAction.Run(upgradeParams.ReleaseName); err == driver.ErrReleaseNotFound {
            logging.Log().Infof("Release %s in namespace %s not found, will run helm install command",
                upgradeParams.ReleaseName, upgradeParams.Namespace)
            return ex.helmInstall(&upgradeParams.InstallParams)
        } else if err != nil {
            return err
        }
    }

    chartName := getChartName(&upgradeParams.InstallParams)
    chartPath, err := upgradeAction.ChartPathOptions.LocateChart(chartName, helmSettings)
    if err != nil {
        logging.Log().Errorf("An error occurred while downloading helm chart: %s. Error message: %s", chartName, err.Error())
        return errors.NewHelmError(err.Error())
    }

    chart, err := getChart(chartPath)
    if err != nil {
        return err
    }

    mergedValues, err := ex.mergeValues(helmSettings, &upgradeParams.InstallParams)
    if err != nil {
        return err
    }

    logging.Log().Infof("Starting helm upgrade command for release %s in namespace %s",
        upgradeParams.ReleaseName, upgradeParams.Namespace)
    _, err = upgradeAction.Run(upgradeParams.ReleaseName, chart, mergedValues)
    if err != nil {
        logging.Log().Errorf("An error occurred while running helm upgrade command. Error message: %s", err.Error())
        return errors.NewHelmError(err.Error())
    }
    return nil
}

func (ex *executor) configureUpgradeAction(helmSettings *cli.EnvSettings, upgradeParams *model.UpgradeParams) (*action.Upgrade, error) {
    actionConfig, err := ex.getActionConfig(helmSettings)
    if err != nil {
        return nil, err
    }

    upgradeAction := action.NewUpgrade(actionConfig)
    upgradeAction.Namespace = upgradeParams.Namespace
    upgradeAction.Timeout = time.Second * time.Duration(upgradeParams.Timeout)
    upgradeAction.MaxHistory = upgradeParams.MaxHistory
    upgradeAction.DisableOpenAPIValidation = upgradeParams.DisableOpenAPIValidation
    upgradeAction.Wait = upgradeParams.Wait
    upgradeAction.RepoURL = upgradeParams.ChartRepo
    upgradeAction.DisableHooks = upgradeParams.DisableHooks
    upgradeAction.Atomic = upgradeParams.Atomic
    upgradeAction.Install = upgradeParams.Install

    return upgradeAction, nil
}

func (ex *executor) configureHistoryAction(helmSettngs *cli.EnvSettings, baseParams *model.BaseParams) (*action.History, error) {
    actionConfig, err := ex.getActionConfig(helmSettngs)
    if err != nil {
        return nil, err
    }

    historyAction := action.NewHistory(actionConfig)
    historyAction.Max = 1
    return historyAction, nil
}
