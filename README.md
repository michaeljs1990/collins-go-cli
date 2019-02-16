Collins Golang CLI
==================

Command line interface for interacting with https://tumblr.github.io/collins/

[![asciicast](https://asciinema.org/a/HfytRKrk8jpgVFmFxOiZyusFS.svg)](https://asciinema.org/a/HfytRKrk8jpgVFmFxOiZyusFS)

## Config file

In order to talk to collins you will need to create a `~/.collins.yml` file with the following fields.

```
host: http://10.0.0.5:9000
username: blake
password: test
```

## Building

You can easily build this by cloning the repo and running `go build -mod=vendor`. Additionally
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
$ bazel build //:collins-go-cli
```

## Coming Features

* Collins IPAM subcommand

## Documentation

All docs can be found in markdown under `docs/`.
