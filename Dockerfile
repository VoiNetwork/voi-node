ARG BASE_ALGORAND_VERSION="3.23.1"
ARG TARGETPLATFORM
FROM algorand/stable:${BASE_ALGORAND_VERSION} AS algorand

FROM --platform=$BUILDPLATFORM golang:1.22 AS builder
WORKDIR /
COPY ./tools/ /tools
COPY Makefile /
COPY go.mod /

RUN echo "Building for $TARGETPLATFORM"
RUN GOOS=linux GOARCH=$(echo $TARGETPLATFORM | cut -d'/' -f2) make all

FROM gcr.io/distroless/cc AS distroless
ENV TELEMETRY_NAME="${HOSTNAME}"
ENV VOINETWORK_PROFILE="${VOINETWORK_PROFILE}"
ENV VOINETWORK_NETWORK="${VOINETWORK_NETWORK}"
ENV VOINETWORK_CATCHUP="${VOINETWORK_CATCHUP}"
ENV VOINETWORK_GENESIS="${VOINETWORK_GENESIS}"
ENV VOINETWORK_CONFIGURATION="${VOINETWORK_CONFIGURATION}"
ENV VOINETWORK_INCOMING_CONNECTIONS_LIMIT="${VOINETWORK_INCOMING_CONNECTIONS_LIMIT}"

HEALTHCHECK --interval=5s --timeout=10s --retries=3 --start-period=10s CMD ["/node/bin/algodhealth"]

COPY --from=algorand --chown=0:0 /root/node/algod /node/bin/algod
COPY --from=algorand --chown=0:0 /root/node/goal /node/bin/goal
COPY --from=algorand --chown=0:0 /root/node/node_exporter /node/bin/node_exporter
COPY --from=builder /build/ /node/bin/
COPY configuration /algod/configuration

CMD ["/node/bin/start-node"]
