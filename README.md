# hooktftp

Hooktftp is a dynamic read-only TFTP server. It's dynamic in a sense it is
executes hooks matched on read requests (RRQ) instead of reading files from
the file system. Hooks are matched with regular expressions and on match
hooktftp will execute a script, issue a HTTP GET request or just reads the file
from the filesystem.

This is inspired by [puavo-tftp]. It's written in Go in the hope of being faster
and more stable.

## Usage

    hooktftp [-v] [config]

Config will be read from `/etc/hooktftp.yml` by default. Verbose option `-v`
print log to stderr insteadof syslog.

## Configuration

Configuration file is in yaml format and it can contain following keys:

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

Regexp can be expanded to the template using the `$0`, `$1`, `$2` etc.
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

## Install

### Compiling from sources

[Install Go][] (1.8 or later), make sure you have git and bazaar too.
Assuming you've successfully set up GOPATH and have GOPATH/bin on your path, simply:
    
    make build
    
Now you should have a standalone hooktftp binary on your path.

    hooktftp -h
    Usage: hooktftp [-v] [config]

### Docker

Try hooktftp by using the [official Docker image](https://hub.docker.com/r/tftpgoteam/hooktftp/):

    $> docker pull tftpgoteam/hooktftp
    $> docker run --rm -ti -v /tmp/myfiles:/var/lib/tftpboot tftpgoteam/hooktftp

### Build Debian package

The package has been created with devscripts and dh-make. To build it:

    debuild -e GOPATH=$PWD -us -uc

## History

Please read debian changelog file.

## Tests

You can start unit and end-to-end test with Makefile target `test`:

    make test

## Hack with Docker

The easiest way to start hacking with hooktftp is to use Docker.

```
# Clone the repository
[host]> git clone git@github.com:tftp-go-team/hooktftp.git
[host]> cd hooktftp

# Create the configuration file. Edit this file later to change hooktftp configuration.
[host]> touch /tmp/hooktftp.yml

# Run the official golang image
[host]> docker run --rm -ti \
    --name tftp-server \
    -v /tmp/hooktftp.yml:/etc/hooktftp.yml:ro \
    -v `pwd`:/go/src/github.com/tftp-go-team/hooktftp \
    golang

# From the container, run hooktftp
[container]> cd src/github.com/tftp-go-team/hooktftp/
[container]> make
[container]> ./src/hooktftp -v

# To query the TFTP server, run the client on another container
[host]> docker run --rm -ti --link tftp-server ubuntu bash
[container]> apt-get update && apt-get install -y tftp-hpa
[container]> echo binary $'\n' get myfile | tftp tftp-server
```

[epeli/hooktftp]: https://github.com/epeli/hooktftp
[puavo-tftp]: https://github.com/opinsys/puavo-tftp
[Install Go]: http://golang.org/doc/install
