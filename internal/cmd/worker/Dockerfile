
ARG BUILD_BASE=cgr.dev/chainguard/go
ARG FINAL_BASE=cgr.dev/chainguard/alpine-base

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
RUN echo -e "\nhttp://dl-cdn.alpinelinux.org/alpine/v3.18/community" >> /etc/apk/repositories
RUN apk add --update bash curl jq cosign

# anchore tools (grype)
RUN curl -sSfL https://raw.githubusercontent.com/anchore/grype/main/install.sh \
    | sh -s -- -b /usr/local/bin

# aquasecurity tools (trivy)
RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh \
    | sh -s -- -b /usr/local/bin

# snyk tools (snyk)
RUN curl -sfL https://static.snyk.io/cli/latest/snyk-linux \
    -o /usr/local/bin/snyk && chmod +x /usr/local/bin/snyk

# cleanup
RUN rm -rf /var/cache/apk/*

COPY --from=builder /src/worker /worker
ENTRYPOINT ["/worker"]
