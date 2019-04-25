#!/bin/bash
# From everything I have found it doesn't look like you can do a bazel build with
# multiple different platforms in one invocation. If that changes we can remove
# this but for now to cut a release you can just run this script.

builds=(
"solaris_amd64"
"openbsd_amd64"
"netbsd_amd64"
"linux_amd64"
"freebsd_amd64"
"darwin_amd64")

for arch in "${builds[@]}"; do
  echo build for $arch
  echo ===========================================================================
  bazel build --platforms=@io_bazel_rules_go//go/toolchain:$arch //cmd/collins:collins
done

releases="release/$(git describe)"
mkdir -p ${releases}

for arch in "${builds[@]}"; do
  cp -f $(pwd)/bazel-bin/cmd/collins/${arch}_pure_stripped/collins $(pwd)/${releases}/collins_$arch
done
