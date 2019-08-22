FROM golang:latest as build
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.3/dep-linux-amd64 && chmod +x /usr/local/bin/dep
COPY Gopkg.toml Gopkg.lock .ssh/id_rsa /go/src/stash.tutu.ru/opscore-workshop-admin/request-generator/
WORKDIR /go/src/stash.tutu.ru/opscore-workshop-admin/request-generator/
RUN go version
RUN mkdir -p ~/.ssh && cp id_rsa ~/.ssh/id_rsa && \
    git config --global url."ssh://git@depot.tutu.ru:7999/".insteadOf "https://stash.tutu.ru/scm/" && \
    git config --global core.sshCommand "ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no" && \
    dep ensure -vendor-only && \
    rm -r id_rsa ~/.ssh
COPY . /go/src/stash.tutu.ru/opscore-workshop-admin/request-generator/
RUN export BUILD_VERSION="$(git name-rev HEAD --name-only)@$(git rev-parse --short HEAD)" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build \
            -ldflags="-X opscore-workshop-admin/request-generator/vendor/stash.tutu.ru/golang/resources/sentry.Version=$BUILD_VERSION" \
            -o /bin/request-generator ./cmd/

FROM alpine:latest as production
COPY --from=build /bin/request-generator /app/request-generator
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY ./openapi.json /openapi.json
RUN apk --update add ca-certificates
ENV TZ "Europe/Moscow"
ENV PORT 7784
EXPOSE $PORT
ENTRYPOINT /app/request-generator
