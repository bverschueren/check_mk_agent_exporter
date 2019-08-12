FROM centos:7

ARG PORT=2112
ARG OS=linux
ARG ARCH=amd64
ARG UID=1000
ARG GID=1000

RUN mkdir /app && \
    chown -R ${UID}:${GID} /app && \
    groupadd --gid ${GID} app && \
    useradd --no-create-home --uid ${UID} --gid ${GID} --home-dir /app app


COPY bin/${OS}_${ARCH}/check_mk_exporter  /bin/check_mk_exporter

USER ${UID}
VOLUME /app/.ssh/
WORKDIR /app

EXPOSE $PORT

CMD  [ "/bin/check_mk_exporter" ]
