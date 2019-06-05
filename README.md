# check_mk_agent exporter

Prometheus exporter collecting metrics from check_mk_agent directly over SSH.

## Usage

```sh
./check_mk_exporter --help
usage: check_mk_exporter [<flags>]

Flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
      --config.file="/etc/check_mk_exporter/ssh.yaml"
                             Config file to use
      --listen.port=2112     Port to listen on
  -l, --log.level=LOG.LEVEL  Enable specify log level

```
Specify the target in the GET call:
```sh
curl "http://localhost:2112/check_mk?target=myhost01
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

These properties can be overruled using query parameters:

 ```sh
curl "http://localhost:2112/check_mk?target=myhost01&port=2222"
```

## Collectors

Currently included collectors:

 - df
 - diskstat

