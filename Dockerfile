# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-tatari-dev-db"
LABEL REPO="https://github.com/mikel-at-tatari/tatari-dev-db"

ENV PROJPATH=/go/src/github.com/mikel-at-tatari/tatari-dev-db

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/mikel-at-tatari/tatari-dev-db
WORKDIR /go/src/github.com/mikel-at-tatari/tatari-dev-db

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/mikel-at-tatari/tatari-dev-db"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/tatari-dev-db/bin

WORKDIR /opt/tatari-dev-db/bin

COPY --from=build-stage /go/src/github.com/mikel-at-tatari/tatari-dev-db/bin/tatari-dev-db /opt/tatari-dev-db/bin/
RUN chmod +x /opt/tatari-dev-db/bin/tatari-dev-db

# Create appuser
RUN adduser -D -g '' tatari-dev-db
USER tatari-dev-db

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/tatari-dev-db/bin/tatari-dev-db"]
