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
package model

type CommandContext struct {
    Version           string                 `json:"version" validate:"required,eq=v1"`
    HelmClientVersion string                 `json:"helmClientVersion" validate:"required"`
    CommandType       string                 `json:"commandType" validate:"required"`
    CommandParams     map[string]interface{} `json:"commandParams" validate:"required"`
}

type BaseParams struct {
    Namespace            string `mapstructure:"namespace" validate:"required"`
    ReleaseName          string `mapstructure:"releaseName" validate:"required"`
    ClusterConfigFileKey string `mapstructure:"clusterConfigFileContentKey" validate:"required"`
    Timeout              int64  `mapstructure:"timeout" validate:"required"`
    HelmDebug            bool   `mapstructure:"helmDebug"`
}

type InstallParams struct {
    BaseParams               `mapstructure:",squash" validate:"required,dive,required"`
    ValuesFileKey            string   `mapstructure:"valuesFileContentKey"`
    DayZeroValuesFileKey     string   `mapstructure:"additionalValuesFileContentKey"`
    ChartUrl                 string   `mapstructure:"chartUrl" validate:"required_without=ChartName"`
    ChartName                string   `mapstructure:"chartName" validate:"required_without=ChartUrl"`
    ChartVersion             string   `mapstructure:"chartVersion" validate:"required_with=ChartName"`
    ChartRepo                string   `mapstructure:"chartVersion" validate:"required_with=ChartName"`
    SetFlagValues            []string `mapstructure:"setFlagValues"`
    DisableOpenAPIValidation bool     `mapstructure:"disableOpenAPIValidation"`
    Wait                     bool     `mapstructure:"helmWait"`
    DisableHooks             bool     `mapstructure:"helmNoHooks"`
    Atomic                   bool     `mapstructure:"atomic"`
    CreateNamespace          bool     `mapstructure:"createNamespace"`
}

type UpgradeParams struct {
    InstallParams `mapstructure:",squash" validate:"required,dive,required"`
    MaxHistory    int  `mapstructure:"maxHistory"`
    Install       bool `mapstructure:"install"`
}

type RollbackParams struct {
    BaseParams `mapstructure:",squash" validate:"required,dive,required"`
    Revision   string `mapstructure:"revisionNumber" validate:"required"`
    MaxHistory int    `mapstructure:"maxHistory"`
}

type UninstallParams struct {
    BaseParams `mapstructure:",squash" validate:"required,dive,required"`
}
