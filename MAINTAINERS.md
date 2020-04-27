**This guide is intended for hooktftp maintainers. If you are not a maintainer, you probably want to check out the [documentation](README.md) instead.**

## Package release HOWTO

You made some updates on hooktftp and want to release a new version for your users? Make sure to complete this todo list.


### Make sure debian package is still working

The debian/ directory is used to create a .deb package. Update the changelog, and test the package is still working with the following commands:

    $> make shell
    #> apt-get install -y build-essential debhelper golang-go
    #> dpkg-buildpackage -us -uc
    #> cd ..
    #> dpkg -i hooktftp_.deb

### Docker image

Build and release the [Docker image](https://hub.docker.com/r/tftpgoteam/hooktftp/):

    $> make release-docker-image

Ping @brmzkw on Github or send him an email at castets.j - at - gmail.com to ask him to make to release. If you want to do it by yourself, ask him to grant you the permissions to do so.
