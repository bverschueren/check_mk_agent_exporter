FROM centos:7

ARG UID=1000
ARG GID=1000
ENV USER=app

RUN yum update -y \
	&& yum install -y openssh-server \
	&& yum clean all \
	&& ssh-keygen -A \
	&& sed -i 's/#PermitRootLogin/PermitRootLogin/' /etc/ssh/sshd_config

RUN groupadd --gid ${GID} ${USER} \
        && useradd --home /app --uid ${UID} --gid ${GID} ${USER} \
	&& echo "cat /usr/local/share/check_mk.txt" > /usr/local/bin/check_mk_agent \
	&& chmod +x /usr/local/bin/check_mk_agent

EXPOSE 22
VOLUME /app/.ssh/
WORKDIR /app

CMD ["/usr/sbin/sshd", "-D", "-e", "-d"]
