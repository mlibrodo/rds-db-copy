# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-rds-db-copy"
LABEL REPO="https://github.com/mlibrodo/rds-db-copy"

ENV PROJPATH=/go/src/github.com/mlibrodo/rds-db-copy

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/mlibrodo/rds-db-copy
WORKDIR /go/src/github.com/mlibrodo/rds-db-copy

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/mlibrodo/rds-db-copy"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/rds-db-copy/bin

WORKDIR /opt/rds-db-copy/bin

COPY --from=build-stage /go/src/github.com/mlibrodo/rds-db-copy/bin/rds-db-copy /opt/rds-db-copy/bin/
RUN chmod +x /opt/rds-db-copy/bin/rds-db-copy

# Create appuser
RUN adduser -D -g '' rds-db-copy
USER rds-db-copy

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/rds-db-copy/bin/rds-db-copy"]
