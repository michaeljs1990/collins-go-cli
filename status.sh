#!/bin/bash
# Produce environment variable that we can use inside out build files.
# Mostly used to pass in environment variables to golang ldflags.
# Vars with STABLE_* format cause builds to retrigger.
echo "STABLE_GIT_COMMIT $(git rev-parse HEAD)"
echo "STABLE_GIT_TAG $(git describe)"
