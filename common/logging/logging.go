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
package logging

import (
    "gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common/config"
    "github.com/openzipkin/zipkin-go"
    "github.com/openzipkin/zipkin-go/reporter"
    "github.com/sirupsen/logrus"
    lSyslog "github.com/sirupsen/logrus/hooks/syslog"
    "log/syslog"
    "net"
    "strconv"
)

const serviceId = "eric-eo-vnfm-helm-executor"

type Logger struct {
    *logrus.Entry
}

var appLogger *Logger

func InitLogger(config *config.Config) {
    fields := buildFields(config.TraceId)
    logger := logrus.New()
    logger.SetReportCaller(true)
    logger.SetFormatter(&logrus.JSONFormatter{
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyLevel: "severity",
            logrus.FieldKeyTime:  "timestamp",
            logrus.FieldKeyMsg:   "message",
            logrus.FieldKeyFunc:  "logger",
        },
        DisableHTMLEscape: true,
    })
    debug, _ := strconv.ParseBool(config.HelmDebug)
    if debug {
        logger.Level = logrus.DebugLevel
    }

    if syslogHook := buildSyslogHook(config, debug); syslogHook != nil {
        logger.AddHook(syslogHook)
    } else {
        logger.Warning("Unable to set syslog hook, logs won't be streamed to log transformer")
    }

    entry := logger.WithFields(fields)
    appLogger = &Logger{entry}
}

func Log() *Logger {
    return appLogger
}

func buildFields(traceId string) logrus.Fields {
    fields := logrus.Fields{
        "service_id": serviceId,
        "traceId":    traceId,
    }

    spanId := getSpanID()
    if spanId != "" {
        fields["spanId"] = spanId
    }

    return fields
}

func getSpanID() string {
    tracer, err := zipkin.NewTracer(reporter.NewNoopReporter())
    if err != nil {
        return ""
    }
    span := tracer.StartSpan(serviceId)
    return span.Context().ID.String()
}

func buildSyslogHook(config *config.Config, debug bool) *lSyslog.SyslogHook {
    if config.LogstashHost == "" || config.LogstashPort == "" {
        return nil
    }
    logStashHostPort := net.JoinHostPort(config.LogstashHost, config.LogstashPort)
    logLevel := syslog.LOG_INFO
    if debug {
        logLevel = syslog.LOG_DEBUG
    }
    syslogHook, err := lSyslog.NewSyslogHook("tcp", logStashHostPort, logLevel, "")
    if err != nil {
        return nil
    }
    return syslogHook
}
