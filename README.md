# docker-machine-driver-hyperkit

The Hyperkit driver will eventually replace the existing xhyve driver and uses [moby/hyperkit](http://github.com/moby/hyperkit) as a Go library.

To install the hyperkit driver:

```shell
make build
```

The hyperkit driver currently requires running as root to use the vmnet framework to setup networking.

If you encountered errors like `Could not find hyperkit executable`, you might need to install [Docker for Mac](https://store.docker.com/editions/community/docker-ce-desktop-mac)
