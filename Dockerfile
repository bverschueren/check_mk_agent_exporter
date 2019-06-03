FROM        centos:7

COPY bin/check_mk_exporter  /bin/check_mk_exporter

EXPOSE      2112
ENTRYPOINT  [ "/bin/check_mk_exporter" ]
