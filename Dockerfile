ARG BASE_ALGORAND_VERSION="3.25.0-stable"
FROM algorand/algod:$BASE_ALGORAND_VERSION AS algod
FROM urtho/algod-voitest-rly:latest AS urtho

FROM golang:1.22 AS builder
WORKDIR /src
COPY ./tools/ /src/
RUN CGO_ENABLED=0 go build -o /dist/algodhealth ./algodhealth.go && \
    CGO_ENABLED=0 go build -o /dist/catch-catchpoint ./catch-catchpoint.go && \
    CGO_ENABLED=0 go build -o /dist/start-node ./start-node.go && \
    CGO_ENABLED=0 go build -o /dist/get-metrics ./get-metrics.go && \
    CGO_ENABLED=0 go build -o /dist/start-metrics ./start-metrics.go

FROM gcr.io/distroless/cc as distroless
ENV TELEMETRY_NAME="${HOSTNAME}"

HEALTHCHECK --interval=5s --timeout=10s --retries=3 --start-period=10s CMD ["/node/bin/algodhealth"]

COPY --from=algod --chown=0:0 /node/bin/algod /node/bin/algod
COPY --from=algod --chown=0:0 /node/bin/goal /node/bin/goal
COPY --from=urtho --chown=0:0 /node/node_exporter /node/bin/node_exporter
COPY --from=builder /dist/ /node/bin/
COPY configuration /algod/configuration

CMD ["/node/bin/start-node"]
