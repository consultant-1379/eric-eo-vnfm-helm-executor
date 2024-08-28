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
package message

import (
    "fmt"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/model"
)

const (
    integerType       = "Integer"
    stringType        = "String"
    commandOutput     = "commandOutput"
    commandExitStatus = "commandExitStatus"
    messageName       = "HelmCommandLifeCycleNotification"
    successExitValue  = "0"
    failureExitValue  = "1"
)

type CamundaMessage struct {
    ProcessInstanceId string                 `json:"processInstanceId"`
    ProcessVariables  map[string]interface{} `json:"processVariables"`
    MessageName       string                 `json:"messageName"`
    ResultEnabled     bool                   `json:"resultEnabled"`
}

func BuildCamundaMessage(err error, processInstanceId string, commandCtx *model.CommandContext) *CamundaMessage {
    processVariables := buildProcessVariables(err, commandCtx)

    return &CamundaMessage{
        ProcessInstanceId: processInstanceId,
        ProcessVariables:  processVariables,
        MessageName:       messageName,
        ResultEnabled:     true,
    }
}

func buildProcessVariables(err error, commandCtx *model.CommandContext) map[string]interface{} {
    resultMsg, exitValue := buildCommandOutput(err, commandCtx)

    processVariables := make(map[string]interface{})
    processVariables[commandOutput] = buildValue(resultMsg, stringType)
    processVariables[commandExitStatus] = buildValue(exitValue, integerType)

    return processVariables
}

func buildCommandOutput(err error, commandCtx *model.CommandContext) (resultMsg string, exitValue string) {
    if err != nil {
        resultMsg = err.Error()
        exitValue = failureExitValue
        return
    }

    namespace := commandCtx.CommandParams["namespace"]
    releaseName := commandCtx.CommandParams["releaseName"]
    resultMsg = fmt.Sprintf(
        "Helm %s command for %s release in %s namespace completed successfully",
        commandCtx.CommandType,
        releaseName,
        namespace,
    )
    exitValue = successExitValue
    return
}

func buildValue(value, valueType string) (m map[string]string) {
    m = make(map[string]string)
    m["type"] = valueType
    m["value"] = value
    return
}
