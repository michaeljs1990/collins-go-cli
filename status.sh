#!/bin/bash
# Produce environment variable that we can use inside our build files.
# Mostly used to pass in environment variables to golang ldflags.
# Vars with STABLE_* format cause builds to retrigger when evaluated
# by bazel.
echo "STABLE_GIT_COMMIT $(git rev-parse HEAD)"
echo "STABLE_GIT_TAG $(git describe)"
