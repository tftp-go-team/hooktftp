[![Build Status](https://travis-ci.org/epeli/hooktftp.png?branch=master)](https://travis-ci.org/epeli/hooktftp)

# hooktftp

Hook based tftp server written in Go.

## Usage

    hooktftp [config]

Config will be read from `/etc/hooktftp.yml` by default.

## Config

Config file is in yaml format and it can contain following keys:

  - `port`: Port to listen (required).
  - `user`: User to drop privileges to.
  - `hooks`: Array of hooks. One or more is required.

### Hooks

Hooks consists of following keys:

  - `name`: Name of the hook.
  - `regexp`: Regexp matcher.
    - Hook is executed when this regexp matches the requested path.
  - `file`, `url` or `shell`: A template where the regexp is expanded.
    - One is required.
    - Determines the type of the hook

Regexp can be expanded to the template using the `$0`, `$1`, `$1` etc.
variables. `$0` is the full matched regexp and rest are the matched regexp
groups.

For example to share files from `/var/lib/tftpboot` add following hook:

```yaml
name: Boot files
regexp: ^.*$
file: /var/lib/tftpboot/$0
```

Or to share custom boot configurations for PXELINUX from a http server based on
mac address:

```yaml
name: Boot configurations
regexp: pxelinux.cfg\/01-(([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2})
url: http://localhost:8080/bootconfigurations/$1
```

The order of the hooks matter. The first one matched is used.

To put it all together:


```yaml
port: 1234
user: hooktftp
hooks:
  - name: Boot files
    regexp: ^.*$
    file: /var/lib/tftpboot/$0

  - name: Boot configurations
    regexp: pxelinux.cfg\/01-(([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2})
    url: http://localhost:8080/bootconfigurations/$1
```

# Releases

See <https://github.com/epeli/hooktftp/releases>
