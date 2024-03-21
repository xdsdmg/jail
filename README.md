# Simple Container

Since my work of nearly half a year is related to containers, I implemented a simple container  just out of curiosity.

**This is only a toy, not use it in the formal environment**.

Implements a simple container model through [Linux Namespace](https://lwn.net/Articles/531114/) and [UnionFS](https://en.wikipedia.org/wiki/UnionFS), and each container has isolated filesystem isolation from network systems.

## Reference

- https://coolshell.cn/articles/17010.html

## Requirement

- **Only test in Debian**.
- Golang 1.19

## Usage

Create RootFS, RootFS will be used by container to generate it's own file system.

``` bash
make init
```

Build server and container, container will make a request to server for network configuration through [UNIX domain socket](https://en.wikipedia.org/wiki/Unix_domain_socket).

``` bash
make build
```

Run server, **server needs to be started before container**.

```
make run_server
```

Create a container, you can create multiple containers, and the network between the containers is unimpeded.

```
make run_container
```

Clean the environment.

```
make clean
```

## TODO

- `umount` UnionFS may fail when quit container.