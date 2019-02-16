Collins Golang CLI
==================

Command line interface for interacting with https://tumblr.github.io/collins/

[![asciicast](https://asciinema.org/a/HfytRKrk8jpgVFmFxOiZyusFS.svg)](https://asciinema.org/a/HfytRKrk8jpgVFmFxOiZyusFS)

## Installing

You can easily build this by cloning the repo and running `go build -mod=vendor`. Additionally
the most recent tags will have binaries uploaded for most common and some uncommon platforms
that you can pull down via curl/wget.

```bash
# Linux Easy Install
$ curl -o collins -L https://github.com/michaeljs1990/collins-go-cli/releases/download/0.9.0/collins-go-cli_linux_amd64
$ chmod +x collins
```

## Building

This repos supports building with bazel as well as the default `go build` command. If you would like to build with bazel
which I use for releases now you can look at github.com/bazelbuild/bazel-gazelle/cmd/gazelle which is used to generate the
build files. However it's not 100% feature complete and running it in this repo without knowing a little about bazel will
likely cause some breakage. If you would like to build the world you can run the following which will pull in some git
keys that are created with the `setup.sh` script and then build the binary.

```
$ bazel build //...
```

## Coming Features

* Collins IPAM subcommand

## Documentation

All docs can be found in markdown under `docs/`.
