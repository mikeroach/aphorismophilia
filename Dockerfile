# Use multi-stage builds per https://docs.docker.com/develop/develop-images/multistage-build/
# to facilitate local development and testing in a build-like environment, but complicates
# Jenkins pipeline due to https://issues.jenkins-ci.org/browse/JENKINS-44609 . See companion
# Jenkinsfile in this repository for more information.

# Per https://github.com/docker-library/docs/blob/f082a22d9ff958fd91d24b94c0e5d4a0af69cf1d/golang/variant-alpine.md -
# consider migrating to a supported golang builder image variant. Look into preserving
# small release artifact sizes with a distroless image (account for fortune-mod).

# Stage 0 (Alpine base image + Golang tools for build & test dependencies)
FROM golang:1.19-alpine3.17 as base
# FIXME: This fortune package only returns offensive fortunes, breaking "mode" option compatibility for the fortune module.
WORKDIR /tmp
COPY vendor ./vendor
RUN apk add --repositories-file=/dev/null --allow-untrusted --no-network --no-cache vendor/fortune-0.1-r1.apk vendor/libbsd-0.9.1-r0.apk
RUN tar xf vendor/gotestsum_0.3.5_linux_amd64.tar.gz && cp gotestsum /usr/bin/

# Stage 1 (test)
FROM base as test
WORKDIR /go/src/aphorismophilia
COPY . .
#RUN CGO_ENABLED=0 go test -v ./...
#RUN echo "Skip unit tests and just test the Docker build."
# FIXME: Exit cleanly regardless of exit code so we can import junit results into Jenkins.
RUN CGO_ENABLED=0 gotestsum --junitfile ut-results.xml -f standard-verbose ; exit 0

# Stage 2 (build and install)
FROM base as build
ARG BUILD=docker
WORKDIR /go/src/aphorismophilia
COPY . .
RUN go install -ldflags "-s -w -X main.build=${BUILD}" -v ./...
RUN cp -a /go/src/aphorismophilia/backends/flatfile/wisdom.txt /go/bin/

# Stage 3 (release)
# FIXME: This uses a separate minimal Alpine image without the Go compiler and test tool
# included from the Golang-built container above. This is cheating since we're not deploying
# the same environment we just tested - to keep things honest we should probably start with
# the same Alpine base image and manage our own build/test containers (e.g. duplicate Golang's
# Dockerfile) to ensure maximum environmental consistency.
FROM alpine:3.17 as release
# FIXME: This package only returns offensive fortunes, breaking "mode" option compatibility for the fortune module.
WORKDIR /tmp
COPY vendor/*.apk ./
RUN apk add --repositories-file=/dev/null --allow-untrusted --no-network --no-cache ./fortune-0.1-r1.apk ./libbsd-0.9.1-r0.apk ; rm ./fortune-0.1-r1.apk ./libbsd-0.9.1-r0.apk
WORKDIR /go/bin
COPY --from=build /go/bin/* ./

EXPOSE 8888
CMD ["/go/bin/aphorismophilia"]