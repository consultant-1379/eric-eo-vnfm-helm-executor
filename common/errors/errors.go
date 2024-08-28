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
package errors

import "errors"

type HelmError struct {
    message string
}

func NewHelmError(message string) *HelmError {
    return &HelmError{
        message: message,
    }
}

func (err *HelmError) Error() string {
    return err.message
}

func GetExitCode(err error) int {
    var helmError *HelmError
    if err == nil || errors.As(err, &helmError) {
        return 0
    }
    return 1
}
