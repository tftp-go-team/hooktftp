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

  - `type`: Type of the hook
    - `file`: Get response from file system
    - `http`: Get reponse from a HTTP server
    - `shell`: Get response from external application
  - `regexp`: Regexp matcher.
    - Hook is executed when this regexp matches the requested path.
  - `template`: A template where the regexp is expanded.

Regexp can be expanded to the template using the `$0`, `$1`, `$1` etc.
variables. `$0` is the full matched regexp and rest are the matched regexp
groups.

### Example

Share files from `/var/lib/tftpboot` add following hook:

```yaml
type: file
regexp: ^.*$
template: /var/lib/tftpboot/$0
```

Share custom boot configurations for PXELINUX from a custom http server:

```yaml
type: http
regexp: pxelinux.cfg\/01-(([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2})
template: http: http://localhost:8080/bootconfigurations/$1
```

The order of the hooks matter. The first one matched is used.

To put it all together:


```yaml
port: 69
user: hooktftp
hooks:
  - type: file
    regexp: ^.*$
    template: /var/lib/tftpboot/$0

  - type: http
    regexp: pxelinux.cfg\/01-(([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2})
    template: http://localhost:8080/bootconfigurations/$1
```

# Downloads

See <https://github.com/epeli/hooktftp/releases>
