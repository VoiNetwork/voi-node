ARG BASE_ALGORAND_VERSION="3.25.0-stable"
FROM algorand/algod:$BASE_ALGORAND_VERSION AS algod
FROM prom/node-exporter:latest AS node-exporter

FROM golang:1.22 AS builder
WORKDIR /
COPY ./tools/ /tools
COPY Makefile /
COPY go.mod /
RUN make all

FROM gcr.io/distroless/cc AS distroless
ENV TELEMETRY_NAME="${HOSTNAME}"
ENV VOINETWORK_NETWORK="${VOINETWORK_NETWORK}"
ENV VOINETWORK_CATCHUP="${VOINETWORK_CATCHUP}"
ENV VOINETWORK_GENESIS="${VOINETWORK_GENESIS}"
ENV VOINETWORK_CONFIGURATION="${VOINETWORK_CONFIGURATION}"

HEALTHCHECK --interval=5s --timeout=10s --retries=3 --start-period=10s CMD ["/node/bin/algodhealth"]

COPY --from=algod --chown=0:0 /node/bin/algod /node/bin/algod
COPY --from=algod --chown=0:0 /node/bin/goal /node/bin/goal
COPY --from=node-exporter --chown=0:0 /bin/node_exporter /node/bin/node_exporter
COPY --from=builder /build/ /node/bin/
COPY configuration /algod/configuration

CMD ["/node/bin/start-node"]
