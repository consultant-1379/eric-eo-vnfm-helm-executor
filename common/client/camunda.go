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
package client

import (
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/httputil"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/message"
    "github.com/hashicorp/go-retryablehttp"
)

type Camunda interface {
    SendMessage(*message.CamundaMessage) error
}

type CamundaClient struct {
    retryClient *retryablehttp.Client
    url         string
}

func NewCamundaClient(url string) *CamundaClient {
    client := retryablehttp.NewClient()
    client.Logger = nil
    return &CamundaClient{
        retryClient: client,
        url:         url,
    }
}

func (c *CamundaClient) SendMessage(msg *message.CamundaMessage) (string, error) {
    bodyBytes, err := httputil.SendPostRequest(c.retryClient, c.url, msg)
    if err != nil {
        return "", err
    }
    return string(bodyBytes), nil
}
