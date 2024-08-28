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

func (ex *executor) helmUninstall(uninstallParams *model.UninstallParams) error {
    helmSettings, err := ex.getHelmSettings(&uninstallParams.BaseParams)
    if err != nil {
        return err
    }

    uninstallAction, err := ex.configureUninstallAction(helmSettings, uninstallParams)
    if err != nil {
        return err
    }

    logging.Log().Infof("Starting helm uninstall command for release %s in namespace %s",
        uninstallParams.ReleaseName, uninstallParams.Namespace)
    _, err = uninstallAction.Run(uninstallParams.ReleaseName)
    if err != nil {
        logging.Log().Errorf("An error occurred while running helm uninstall command. Error message: %s", err.Error())
        return errors.NewHelmError(err.Error())
    }
    return nil
}

func (ex *executor) configureUninstallAction(helmSettings *cli.EnvSettings, uninstallParams *model.UninstallParams) (*action.Uninstall, error) {
    actionConfig, err := ex.getActionConfig(helmSettings)
    if err != nil {
        return nil, err
    }

    uninstallAction := action.NewUninstall(actionConfig)
    uninstallAction.Timeout = time.Second * time.Duration(uninstallParams.Timeout)
    return uninstallAction, nil
}
