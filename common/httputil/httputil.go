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
package httputil

import (
    "bytes"
    "encoding/json"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/logging"
    "github.com/hashicorp/go-retryablehttp"
    "io"
    "net/http"
)

func SendPostRequest(client *retryablehttp.Client, url string, data interface{}) ([]byte, error) {
    request, err := buildPostRequest(url, data)
    if err != nil {
        logging.Log().Errorf("An error occured while building post request to %s. Error message: %s", url, err.Error())
        return nil, err
    }

    logging.Log().Infof("Sending post request to: %s", url)
    response, err := client.Do(request)
    if response != nil {
        defer response.Body.Close()
    }
    if err != nil {
        logging.Log().Errorf("An error occurred while sending post request to: %s. Error message: %s", url, err.Error())
        return nil, err
    }

    bodyBytes, err := io.ReadAll(response.Body)
    if err != nil {
        logging.Log().Errorf("An error occurred while reading reponse from: %s. Error message: %s", url, err.Error())
        return nil, err
    }
    return bodyBytes, nil
}

func buildPostRequest(url string, data interface{}) (*retryablehttp.Request, error) {
    dataBytes, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    request, err := retryablehttp.NewRequest(http.MethodPost, url, bytes.NewReader(dataBytes))
    if err != nil {
        return nil, err
    }
    request.Header.Add("Content-Type", "application/json")
    request.Header.Add("Accept", "application/json")
    return request, nil
}
