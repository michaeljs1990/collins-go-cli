<h1 align="center">Collins CLI</h1>

<p align="center">Command line interface for interacting with https://tumblr.github.io/collins/</p>

### Install

You can always install from the releases page on github or if you want you can also stay up to date
using the `go get` method. My server can be a little slow sometimes but it works just fine if you wait.

```
go get -u cgit.xrt0x.com/xrt0x/collins-go-cli/cmd/collins
```

### Config file

In order to talk to collins you will need to create a `~/.collins.yml` file with the following fields.

```
host: http://10.0.0.5:9000
username: blake
password: test
```

**NOTE**: if you do not have the password saved in `~/.collins.yml`, the client will prompt for a password
at runtime. But, this will error and fail, if the collins command is being piped from one collins command to next.

### Documentation

All docs can be found in markdown under `collins/docs/`.

### Development

You can easily build this by cloning the repo and running `go build -mod=vendor ./cmd/collins`. Additionally
the most recent tags will have binaries uploaded for most common and some uncommon platforms
that you can pull down via curl/wget.

This repos supports building with bazel as well as the default `go build` command. If you would like to build with bazel
which I use for releases now you can look at github.com/bazelbuild/bazel-gazelle/cmd/gazelle which is used to generate the
build files. However it's not 100% feature complete and running it in this repo without knowing a little about bazel will
likely cause some breakage. If you would like to build the world you can run the following which will pull in some git
keys that are created with the `setup.sh` script and then build the binary.

To get a similar feature set as I have used gox in the past for you can pass the --platforms flag to bazel build. That
flag can take anything that is output by `bazel query 'kind(platform, @io_bazel_rules_go//go/toolchain:all)'` See the
following page for more info https://github.com/bazelbuild/rules_go/blob/master/go/core.rst#cross-compilation.

```
$ bazel build //cmd/collins:collins
```

Setting up collins itself for testing and populating it is beyond the scope of this readme. If you would like to do that
but don't know how to get started send me a message.
