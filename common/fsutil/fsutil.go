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
package fsutil

import (
    "errors"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/logging"
    "io"
    "os"
)

func SaveTempFile(content, pattern string) (path string, err error) {
    tempFile, err := os.CreateTemp("", pattern)
    if tempFile != nil {
        defer func() {
            closeErr := tempFile.Close()
            err = errors.Join(err, closeErr)
        }()
    }
    if err != nil {
        logging.Log().Errorf("An error occurred while creating temp file for %s", pattern)
        return
    }

    _, err = io.WriteString(tempFile, content)
    if err != nil {
        logging.Log().Errorf("An error occurred while writing to temp file %s", tempFile.Name())
        return
    }
    path = tempFile.Name()
    logging.Log().Infof("Saved %s as temp file: %s", pattern, path)
    return
}
