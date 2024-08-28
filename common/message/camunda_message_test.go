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
    "errors"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/model"
    "reflect"
    "testing"
)

func TestBuildCamundaMessage(t *testing.T) {

    camundaMessageWithError := &CamundaMessage{
        MessageName:       messageName,
        ProcessInstanceId: "b42936ff-f03b-11ed-a9b2-2e2955c6af12",
        ResultEnabled:     true,
        ProcessVariables: map[string]interface{}{
            commandOutput: map[string]string{
                "type":  "String",
                "value": "testing with error",
            },
            commandExitStatus: map[string]string{
                "type":  "Integer",
                "value": "1",
            },
        },
    }

    camundaMessageNoError := &CamundaMessage{
        MessageName:       messageName,
        ProcessInstanceId: "b42936ff-f03b-11ed-a9b2-2e2955c6af12",
        ResultEnabled:     true,
        ProcessVariables: map[string]interface{}{
            commandOutput: map[string]string{
                "type":  "String",
                "value": "Helm install command for test-release-1 release in test-ns namespace completed successfully",
            },
            commandExitStatus: map[string]string{
                "type":  "Integer",
                "value": "0",
            },
        },
    }

    commandCtx := &model.CommandContext{
        Version:           "v1",
        CommandType:       "install",
        HelmClientVersion: "3.8",
        CommandParams: map[string]interface{}{
            "namespace":   "test-ns",
            "releaseName": "test-release-1",
        },
    }

    type args struct {
        err               error
        processInstanceId string
        commandCtx        *model.CommandContext
    }

    argsWithError := args{
        err:               errors.New("testing with error"),
        processInstanceId: "b42936ff-f03b-11ed-a9b2-2e2955c6af12",
        commandCtx:        commandCtx,
    }

    argsNoError := args{
        err:               nil,
        processInstanceId: "b42936ff-f03b-11ed-a9b2-2e2955c6af12",
        commandCtx:        commandCtx,
    }

    tests := []struct {
        name string
        args args
        want *CamundaMessage
    }{
        {name: "test with error", args: argsWithError, want: camundaMessageWithError},
        {name: "test no error", args: argsNoError, want: camundaMessageNoError},
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            if got := BuildCamundaMessage(test.args.err, test.args.processInstanceId, test.args.commandCtx); !reflect.DeepEqual(got, test.want) {
                t.Errorf("BuildCamundaMessage() = %v, want %v", got, test.want)
            }
        })
    }
}
