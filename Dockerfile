FROM golang:1.12

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOCACHE=/tmp \
    STI_SCRIPTS_PATH=/usr/libexec/s2i

LABEL io.openshift.s2i.scripts-url=image://${STI_SCRIPTS_PATH}

COPY ./s2i/bin/ ${STI_SCRIPTS_PATH}

RUN mkdir /build \
    && chmod 0777 /build

WORKDIR /build

USER 1001

CMD ["/usr/libexec/s2i/usage"]
