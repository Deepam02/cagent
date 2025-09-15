# syntax=docker/dockerfile:1

# xx is a helper for cross-compilation
FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.7.0 AS xx

FROM golang:1.25.0-alpine@sha256:f18a072054848d87a8077455f0ac8a25886f2397f88bfdd222d6fafbb5bba440 AS build-agent
RUN apk add --no-cache build-base
WORKDIR /app
COPY . ./
ARG GIT_TAG GIT_COMMIT BUILD_DATE
RUN --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=secret,id=telemetry_api_key \
    --mount=type=secret,id=telemetry_endpoint \
    --mount=type=secret,id=telemetry_header \
    sh -c 'TELEMETRY_API_KEY=$(cat /run/secrets/telemetry_api_key 2>/dev/null || echo "") && TELEMETRY_ENDPOINT=$(cat /run/secrets/telemetry_endpoint 2>/dev/null || echo "") && TELEMETRY_HEADER=$(cat /run/secrets/telemetry_header 2>/dev/null || echo "") && go build -trimpath -ldflags "-s -w -X '"'"'github.com/docker/cagent/cmd/root.Version=$GIT_TAG'"'"' -X '"'"'github.com/docker/cagent/cmd/root.Commit=$GIT_COMMIT'"'"' -X '"'"'github.com/docker/cagent/cmd/root.BuildTime=$BUILD_DATE'"'"' -X '"'"'github.com/docker/cagent/internal/telemetry.TelemetryEndpoint=$TELEMETRY_ENDPOINT'"'"' -X '"'"'github.com/docker/cagent/internal/telemetry.TelemetryAPIKey=$TELEMETRY_API_KEY'"'"' -X '"'"'github.com/docker/cagent/internal/telemetry.TelemetryHeader=$TELEMETRY_HEADER'"'"'" -o /agent .'

FROM --platform=$BUILDPLATFORM golang:1.25.0-alpine3.22 AS builder-base
RUN apk add clang
WORKDIR /src
COPY --from=xx / /
ENV CGO_ENABLED=0
ARG TARGETPLATFORM TARGETOS TARGETARCH
ARG GIT_TAG GIT_COMMIT BUILD_DATE

FROM builder-base AS builder
COPY . ./
RUN --mount=type=cache,target=/root/.cache,id=docker-ai-$TARGETPLATFORM \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=secret,id=telemetry_api_key \
    --mount=type=secret,id=telemetry_endpoint \
    --mount=type=secret,id=telemetry_header <<EOT
    set -ex
    TELEMETRY_API_KEY=$(cat /run/secrets/telemetry_api_key 2>/dev/null || echo "")
    TELEMETRY_ENDPOINT=$(cat /run/secrets/telemetry_endpoint 2>/dev/null || echo "")
    TELEMETRY_HEADER=$(cat /run/secrets/telemetry_header 2>/dev/null || echo "")
    xx-go build -trimpath -ldflags "-s -w -X 'github.com/docker/cagent/cmd/root.Version=$GIT_TAG' -X 'github.com/docker/cagent/cmd/root.Commit=$GIT_COMMIT' -X 'github.com/docker/cagent/cmd/root.BuildTime=$BUILD_DATE' -X 'github.com/docker/cagent/internal/telemetry.TelemetryEndpoint=$TELEMETRY_ENDPOINT' -X 'github.com/docker/cagent/internal/telemetry.TelemetryAPIKey=$TELEMETRY_API_KEY' -X 'github.com/docker/cagent/internal/telemetry.TelemetryHeader=$TELEMETRY_HEADER'" -o /binaries/cagent-$TARGETOS-$TARGETARCH .
    xx-verify --static /binaries/cagent-$TARGETOS-$TARGETARCH
    if [ "$TARGETOS" = "windows" ]; then
      mv /binaries/cagent-$TARGETOS-$TARGETARCH /binaries/cagent-$TARGETOS-$TARGETARCH.exe
    fi
EOT

FROM scratch AS local
ARG TARGETOS TARGETARCH
COPY --from=builder /binaries/cagent-$TARGETOS-$TARGETARCH cagent

FROM scratch AS cross
COPY --from=builder /binaries .

FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1
ARG TARGETOS TARGETARCH
COPY --from=builder /binaries/cagent-$TARGETOS-$TARGETARCH /cagent
RUN mkdir /data
ENTRYPOINT ["/cagent"]
