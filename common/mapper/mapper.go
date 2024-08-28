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
package mapper

import (
    "encoding/json"
    "fmt"
    "github.com/go-playground/validator/v10"
    "github.com/mitchellh/mapstructure"
)

func MapToStruct(src, dst interface{}) error {
    switch src.(type) {
    case string:
        srcString := src.(string)
        err := json.Unmarshal([]byte(srcString), dst)
        if err != nil {
            return err
        }
    case map[string]interface{}:
        srcMap := src.(map[string]interface{})
        err := mapstructure.Decode(srcMap, dst)
        if err != nil {
            return err
        }
    default:
        return fmt.Errorf("unsupported type: %T", src)
    }

    validate := validator.New()
    return validate.Struct(dst)
}
