// ******************************************************************************
// COPYRIGHT Ericsson 2024
//
//
//
// The copyright to the computer program(s) herein is the property of
//
// Ericsson Inc. The programs may be used and/or copied only with written
//
// permission from Ericsson Inc. or in accordance with the terms and
//
// conditions stipulated in the agreement/contract under which the
//
// program(s) have been supplied.
// ******************************************************************************
module gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor

go 1.20

replace gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common => ../common

require gerrit.ericsson.se/OSS/com.ericsson.orchestration.mgmt/helm-executor/common v0.0.0-00010101000000-000000000000

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.1 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
