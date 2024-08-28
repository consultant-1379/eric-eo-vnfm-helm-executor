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

import (
    "errors"
    "testing"
)

func TestGetExitCode(t *testing.T) {
    err := errors.New("should be helm error")

    type args struct {
        err error
    }
    tests := []struct {
        name string
        args args
        want int
    }{
        {"should return 0 when err is nil", args{err: nil}, 0},
        {"should return 0 when err is HelmError type", args{err: NewHelmError(err.Error())}, 0},
        {"should return 1", args{err: errors.New("not helm error")}, 1},
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            if got := GetExitCode(test.args.err); got != test.want {
                t.Errorf("GetExitCode() = %v, want %v", got, test.want)
            }
        })
    }
}
