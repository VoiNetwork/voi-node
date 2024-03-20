ARG BASE_ALGORAND_VERSION="3.22.1-stable"
FROM algorand/algod:$BASE_ALGORAND_VERSION as algod
FROM urtho/algod-voitest-rly:latest as urtho

FROM golang:1.22 as builder
WORKDIR /src
COPY ./tools/ /src/
RUN CGO_ENABLED=0 go build -o /dist/algodhealth ./algodhealth.go && \
    CGO_ENABLED=0 go build -o /dist/catch-catchpoint ./catch-catchpoint.go && \
    CGO_ENABLED=0 go build -o /dist/start-node /src/start-node.go

FROM gcr.io/distroless/cc as distroless
ENV TELEMETRY_NAME="${HOSTNAME}"

HEALTHCHECK --interval=30s --timeout=30s --start-period=20s CMD ["/node/bin/algodhealth"]

COPY --from=algod --chown=0:0 /node/bin/algod /node/bin/algod
COPY --from=algod --chown=0:0 /node/bin/goal /node/bin/goal
COPY --from=urtho --chown=0:0 /node/node_exporter /node/bin/node_exporter
COPY --from=builder /dist/ /node/bin/
COPY configuration /algod/configuration

CMD ["/node/bin/start-node"]
