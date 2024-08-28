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
    "strconv"
    "time"
)

func (ex *executor) helmRollback(rollbackParams *model.RollbackParams) error {
    helmSettings, err := ex.getHelmSettings(&rollbackParams.BaseParams)
    if err != nil {
        return err
    }

    rollbackAction, err := ex.configureRollbackAction(helmSettings, rollbackParams)
    logging.Log().Infof("Starting helm rollback command for release %s in namespace %s",
        rollbackParams.ReleaseName, rollbackParams.Namespace)
    err = rollbackAction.Run(rollbackParams.ReleaseName)
    if err != nil {
        logging.Log().Errorf("An error occurred while running helm rollback command. Error message: %s", err.Error())
        return errors.NewHelmError(err.Error())
    }
    return nil
}

func (ex *executor) configureRollbackAction(helmSettings *cli.EnvSettings, rollbackParams *model.RollbackParams) (*action.Rollback, error) {
    actionConfig, err := ex.getActionConfig(helmSettings)
    if err != nil {
        return nil, err
    }

    rollbackAction := action.NewRollback(actionConfig)
    revision, err := strconv.Atoi(rollbackParams.Revision)
    if err != nil {
        return nil, errors.NewHelmError(err.Error())
    }
    rollbackAction.Version = revision
    rollbackAction.MaxHistory = rollbackParams.MaxHistory
    rollbackAction.Timeout = time.Second * time.Duration(rollbackParams.Timeout)

    return rollbackAction, nil
}
