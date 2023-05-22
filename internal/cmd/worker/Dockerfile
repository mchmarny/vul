
ARG BUILD_BASE=golang:latest
ARG FINAL_BASE=ubuntu:latest

FROM $BUILD_BASE AS builder
WORKDIR /src/
COPY . /src/
ARG VERSION
ENV VERSION=${VERSION}
RUN CGO_ENABLED=0 go build -a -trimpath -ldflags="\
    -w -s -X main.version=$VERSION \
	-extldflags '-static'" \
    -mod vendor -o worker internal/cmd/worker/main.go

# BUILD
FROM $FINAL_BASE
WORKDIR /

# packages
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install curl

# anchore tools (grype)
RUN curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh \
    | sh -s -- -b /usr/local/bin

RUN curl -sSfL https://raw.githubusercontent.com/anchore/grype/main/install.sh \
    | sh -s -- -b /usr/local/bin

# aquasecurity tools (trivy)
RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh \
    | sh -s -- -b /usr/local/bin

# snyk tools (snyk)
RUN curl -sfL https://github.com/snyk/cli/releases/latest/download/snyk-linux \
    -o /usr/local/bin/snyk && chmod +x /usr/local/bin/snyk

COPY --from=builder /src/worker /worker
ENTRYPOINT ["/worker"]
