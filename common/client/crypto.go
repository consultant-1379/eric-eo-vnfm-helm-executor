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
    "encoding/json"
    "fmt"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/httputil"
    "github.com/hashicorp/go-retryablehttp"
)

const decryptUrl = "%s/generic/v1/decryption"

type Decryptor interface {
    Decrypt(*DecryptRequest) (string, error)
}

type DecryptRequest struct {
    Ciphertext string `json:"ciphertext"`
}

type decryptResponse struct {
    Plaintext string `json:"plaintext"`
}

type DecryptClient struct {
    retryClient *retryablehttp.Client
    cryptoHost  string
}

func NewDecryptClient(cryptoHost string) *DecryptClient {
    client := retryablehttp.NewClient()
    client.Logger = nil
    return &DecryptClient{
        retryClient: client,
        cryptoHost:  cryptoHost,
    }
}

func (d *DecryptClient) Decrypt(decryptRequest *DecryptRequest) (string, error) {
    var decryptResp decryptResponse

    bodyBytes, err := httputil.SendPostRequest(d.retryClient, fmt.Sprintf(decryptUrl, d.cryptoHost), decryptRequest)
    if err != nil {
        return "", err
    }
    err = json.Unmarshal(bodyBytes, &decryptResp)
    if err != nil {
        return "", err
    }
    return decryptResp.Plaintext, nil
}
