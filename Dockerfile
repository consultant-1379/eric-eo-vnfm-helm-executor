#
# COPYRIGHT Ericsson 2024
#
#
#
# The copyright to the computer program(s) herein is the property of
#
# Ericsson Inc. The programs may be used and/or copied only with written
#
# permission from Ericsson Inc. or in accordance with the terms and
#
# conditions stipulated in the agreement/contract under which the
#
# program(s) have been supplied.
#

ARG BUILD_IMAGE_TAG="1.20"
ARG BASE_IMAGE_VERSION="5.18.0-14"
# Build Image
FROM armdocker.rnd.ericsson.se/proj-am/sles/sles-golang:${BUILD_IMAGE_TAG} AS build

COPY . /src

RUN cd /src/helm-plugin-3.8.1 && \
    go build -ldflags="-s -w" -trimpath -buildmode=plugin

RUN cd /src/helm-plugin-3.10.1 && \
    go build -ldflags="-s -w" -trimpath -buildmode=plugin

RUN cd /src/helm-plugin-3.12.0 && \
    go build -ldflags="-s -w" -trimpath -buildmode=plugin

RUN cd /src/helm-plugin-3.13.0 && \
    go build -ldflags="-s -w" -trimpath -buildmode=plugin

RUN cd /src/helm-plugin-3.14.2 && \
    go build -ldflags="-s -w" -trimpath -buildmode=plugin

RUN cd /src/helm-executor && \
    go build -ldflags="-s -w" -trimpath -o helm-executor

# Helm-Executor Image
FROM armdocker.rnd.ericsson.se/proj-ldc/common_base_os_release/sles:${BASE_IMAGE_VERSION}

COPY --from=build /src/helm-executor/helm-executor /usr/bin/helm-executor
COPY --from=build /src/helm-plugin-3.8.1/helm-plugin-3.8.1.so /usr/bin/helm-plugin-3.8.1.so
COPY --from=build /src/helm-plugin-3.10.1/helm-plugin-3.10.1.so /usr/bin/helm-plugin-3.10.1.so
COPY --from=build /src/helm-plugin-3.12.0/helm-plugin-3.12.0.so /usr/bin/helm-plugin-3.12.0.so
COPY --from=build /src/helm-plugin-3.13.0/helm-plugin-3.13.0.so /usr/bin/helm-plugin-3.13.0.so
COPY --from=build /src/helm-plugin-3.14.2/helm-plugin-3.14.2.so /usr/bin/helm-plugin-3.14.2.so

ARG HELM_EXECUTOR_GID=155463 \
    HELM_EXECUTOR_UID=155463 \
    HELM_EXECUTOR_DATA_DIR="/helm-executor"

RUN echo "${HELM_EXECUTOR_UID}:x:${HELM_EXECUTOR_UID}:${HELM_EXECUTOR_GID}:wfs-user:${HELM_EXECUTOR_DATA_DIR}:/bin/false" >> /etc/passwd && \
    cat /etc/passwd && \
    sed -i "s|root:/bin/bash|root:/bin/false|g" /etc/passwd && \
    chmod -R g=u /usr/bin/helm-executor && \
    chown -h ${HELM_EXECUTOR_UID}:0 /usr/bin/helm-executor && \
    mkdir -p ${HELM_EXECUTOR_DATA_DIR} && \
    chmod -R g=u ${HELM_EXECUTOR_DATA_DIR} && \
    chown -fR ${HELM_EXECUTOR_UID}:0 ${HELM_EXECUTOR_DATA_DIR}

USER $HELM_EXECUTOR_UID:$HELM_EXECUTOR_GID
