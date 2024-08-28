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
    "fmt"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/client"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/errors"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/logging"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/mapper"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/model"
)

const (
    install   = "install"
    upgrade   = "upgrade"
    crd       = "crd"
    rollback  = "rollback"
    uninstall = "uninstall"
)

func RunHelmCommand(commandCtx *model.CommandContext, decryptClient client.Decryptor, redisClient client.RedisClient) error {
    ex := newExecutor(decryptClient, redisClient)
    logging.Log().Infof("Command type: %s", commandCtx.CommandType)
    switch commandCtx.CommandType {
    case install:
        installParams := &model.InstallParams{}
        err := mapper.MapToStruct(commandCtx.CommandParams, installParams)
        if err != nil {
            logging.Log().Errorf("An error occured while mapping install params. Error message: %s", err.Error())
            return err
        }
        return ex.helmInstall(installParams)
    case crd, upgrade:
        upgradeParams := &model.UpgradeParams{}
        err := mapper.MapToStruct(commandCtx.CommandParams, upgradeParams)
        if err != nil {
            logging.Log().Errorf("An error occured while mapping upgrade params. Error message: %s", err.Error())
            return err
        }
        return ex.helmUpgrade(upgradeParams)
    case rollback:
        rollbackParams := &model.RollbackParams{}
        err := mapper.MapToStruct(commandCtx.CommandParams, rollbackParams)
        if err != nil {
            logging.Log().Errorf("An error occured while mapping rollback params. Error message: %s", err.Error())
            return err
        }
        return ex.helmRollback(rollbackParams)
    case uninstall:
        uninstallParams := &model.UninstallParams{}
        err := mapper.MapToStruct(commandCtx.CommandParams, uninstallParams)
        if err != nil {
            logging.Log().Errorf("An error occured while mapping uninstall params. Error message: %s", err.Error())
            return err
        }
        return ex.helmUninstall(uninstallParams)
    default:
        logging.Log().Errorf("Command type %s is not suported", commandCtx.CommandType)
        return errors.NewHelmError(fmt.Sprintf("unsupported command type: %s", commandCtx.CommandType))
    }
}
