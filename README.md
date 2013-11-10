[![Build Status](https://travis-ci.org/epeli/hooktftp.png?branch=master)](https://travis-ci.org/epeli/hooktftp)

# hooktftp

Hook based tftp server inspired by [puavo-tftp]. It's written in Go in the hope
of being faster and more stable.

It's intented to be used with [PXELINUX][] for dynamic mac address based boots.

## Usage

    hooktftp [config]

Config will be read from `/etc/hooktftp.yml` by default.

## Config

Config file is in yaml format and it can contain following keys:

  - `port`: Port to listen (required)
  - `user`: User to drop privileges to
  - `hooks`: Array of hooks. One or more is required

### Hooks

Hooks consists of following keys:

  - `type`: Type of the hook
    - `file`: Get response from the file system
    - `http`: Get response from a HTTP server
    - `shell`: Get response from external application
  - `regexp`: Regexp matcher
    - Hook is executed when this regexp matches the requested path
  - `template`: A template where the regexp is expanded

Regexp can be expanded to the template using the `$0`, `$1`, `$1` etc.
variables. `$0` is the full matched regexp and rest are the matched regexp
groups.

### Example

To share files from `/var/lib/tftpboot` add following hook:

```yaml
type: file
regexp: ^.*$
template: /var/lib/tftpboot/$0
```

Share custom boot configurations for PXELINUX from a custom http server:

```yaml
type: http
regexp: pxelinux.cfg\/01-(([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2})
template: http://localhost:8080/boot/$1
```

The order of the hooks matter. The first one matched is used.

To put it all together:

```yaml
port: 69
user: hooktftp
hooks:
  - type: http
    regexp: pxelinux.cfg\/01-(([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2})
    template: http://localhost:8080/boot/$1

  - type: file
    regexp: ^.*$
    template: /var/lib/tftpboot/$0
```

## Silly benchmarks

On AMD Phenom II X4 965 and Samsung SSD

    time tools/speed.sh

server            | size | concurrency | count | blocksize | time
----------------- |------|-------------|-------|---------- |-----
puavo-tftp (ruby) | 11M  | 1           | 10    | 512       | 0m27.012s
hooktftp   (Go)   | 11M  | 1           | 10    | 512       | 0m16.126s
tftp-hpa   (C)    | 11M  | 1           | 10    | 512       | 0m14.409s


    time tools/speed-concurrent.sh

server            | size | concurrency | count | blocksize | time
----------------- |------|-------------|-------|---------- |-----
puavo-tftp (ruby) | 11M  | 10          | 10    | 512       | 0m59.869s
hooktftp   (Go)   | 11M  | 10          | 10    | 512       | 0m24.531s
tftp-hpa   (C)    | 11M  | 10          | 10    | 512       | 0m10.326s Broken test?


# Downloads

See <https://github.com/epeli/hooktftp/releases>

[puavo-tftp]: https://github.com/opinsys/puavo-tftp
[PXELINUX]: http://www.syslinux.org/wiki/index.php/PXELINUX
