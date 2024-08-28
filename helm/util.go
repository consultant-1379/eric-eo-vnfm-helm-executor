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
package helm

import (
    "fmt"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/errors"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/logging"
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/model"
    "helm.sh/helm/v3/pkg/action"
    "helm.sh/helm/v3/pkg/chart"
    "helm.sh/helm/v3/pkg/chart/loader"
)

func getChartName(installParams *model.InstallParams) string {
    if installParams.ChartUrl == "" {
        return installParams.ChartName
    }
    return installParams.ChartUrl
}

func getChart(chartPath string) (*chart.Chart, error) {
    logging.Log()
    chart, err := loader.Load(chartPath)
    if err != nil {
        logging.Log().Errorf("An error occurred while loading chart %s. Error message: %s", chart.Name(), err.Error())
        return nil, err
    }

    err = checkIfInstallable(chart)
    if err != nil {
        return nil, err
    }

    err = checkDependencies(chart)
    if err != nil {
        return nil, err
    }

    return chart, nil
}

func checkIfInstallable(chart *chart.Chart) error {
    switch chart.Metadata.Type {
    case "", "application":
        return nil
    default:
        logging.Log().Errorf("Chart %s with type %s is not installable", chart.Name(), chart.Metadata.Type)
        return errors.NewHelmError(fmt.Sprintf("chart %s with type %s is not installable", chart.Name(), chart.Metadata.Type))
    }
}

func checkDependencies(chart *chart.Chart) error {
    if chartDependencies := chart.Metadata.Dependencies; chartDependencies != nil {
        err := action.CheckDependencies(chart, chartDependencies)
        if err != nil {
            logging.Log().Errorf("An error occurred while checking chart %s dependencies. Error message: %s", chart.Name(), err.Error())
            return errors.NewHelmError(err.Error())
        }
    }
    return nil
}
