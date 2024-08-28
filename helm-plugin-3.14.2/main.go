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
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/client"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/model"
    _ "helm.sh/helm/v3/pkg/action"
)

func RunHelmCommand(commandCtx *model.CommandContext, decryptClient client.Decryptor, redisClient client.RedisClient) error {
    return helm.RunHelmCommand(commandCtx, decryptClient, redisClient)
}
