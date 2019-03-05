Collins Golang CLI
==================

Command line interface for interacting with https://tumblr.github.io/collins/

## Whats good about this?

The golang collins cli has a few features not availabe in the ruby client. In no particular order.

1. JSON output format is a valid JSON dump
2. YAML output format is valid YAML in all languages and not limited to ruby
3. You can pipe IPs, tags, hostnames, or anything else to `collins query` to get more info about them
4. Additional fields that are computed client side. See `func fieldToAssetStruct` or the docs for the full list
5. `collins provision` has a wait flag for that pesky rate limiter
6. `collins power` issues commands inside goroutines so bad BMCs don't slow you down
7. If you have been using the collins ruby cli it should function 100% the same for basic workflows
8. `collins query` has negative matching `key:~valueidontwant`
9. `collins query` takes a comma delimited list for all values `-n devnode,aaanode` will return both
10. `collins query` flags all take an option ~ before any value to negate them.

## Config file

In order to talk to collins you will need to create a `~/.collins.yml` file with the following fields.

```
host: http://10.0.0.5:9000
username: blake
password: test
```

| WARNING: unlike the ruby collins-cli gem you may not use symbols and strings interchangeably. If collins-go-cli blows up on the first run make sure you don't have `:host`, `:username`,or `:password` in your config file. Note that changing `:host` to `host` will still work in the ruby client as well so it's not a breaking change.|
| --- |

## Install

You can always install from the releases I have produced if you want however you can also stay up to date
using the `go get` method. My server can be a little slow sometimes but it works just fine if you wait.

```
go get -u cgit.xrt0x.com/xrt0x/collins-go-cli/cmd/collins
```

## Development

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

## Documentation

All docs can be found in markdown under `docs/`.
