**This guide is intended for hooktftp maintainers. If you are not a maintainer,
you probably want to check out the [documentation](README.md) instead.**

## Package release HOWTO

You made some updates on hooktftp and want to release a new version for your
users? Make sure to complete this todo list.


### Docker image

Build the Docker image:

    $> docker build -t tftpgoteam/hooktftp:latest .

A docker image needs to be pushed on the [Docker
hub](https://hub.docker.com/r/tftpgoteam/hooktftp/). Ping @brmzkw on Github or
send him an email at castets.j - at - gmail.com to ask him to make the release.
If you want to do it by yourself, ask him to grant you the permissions to do
so.

Push the image:

    $> docker push tftpgoteam/hooktftp:latest
