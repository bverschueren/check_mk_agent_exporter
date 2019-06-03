# check_mk_agent exporter

Prometheus exporter collecting metrics from check_mk_agent directly over SSH.

## Usage

```sh
./check_mk_exporter
```

## SSH configuration

The exporter reads SSH configuration from a `/etc/check_mk_exporter/ssh.yaml` file by default.

Example config:
```YAML
targets:
  myhost01
    HostName: myhost01.my.domain
    User: myuser
    IdentityFile: /home/myuser/.ssh/private_key
  myhost02
    HostName: myhost02.my.domain
    Port: 2222
    IdentityFile: /home/myuser/.ssh/private_key
```
## Collectors

Currently included collectors:

 - df
 - diskstat
