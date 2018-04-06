#!/bin/sh

# gitlab nonsense
[ -z "${BUILD_NUMBER}" ] && export BUILD_NUMBER=${CI_PIPELINE_ID}
[ -z "${PROJECT_NAME}" ] && export PROJECT_NAME=${CI_PROJECT_NAME}
[ -z "${COMMIT_ID}" ] && export COMMIT_ID=${CI_COMMIT_SHA}
[ -z "${GIT_BRANCH}" ] && export GIT_BRANCH=${CI_COMMIT_REF_NAME}

# we're not on a build server..
if [ -z "${BUILD_NUMBER}" ]; then
    echo no build number - not submitting
    exit 0
fi

# this is only a library. for now, we simply tar it up
rm -rf dist ; mkdir dist || exit 10
tar -jcvf dist/go-framework.tar.bz2 -C src go-framework || exit 10


# we are on a build server, so submit it:
build-repo-client -branch=${GIT_BRANCH} -build=${BUILD_NUMBER} -commitid=${COMMIT_ID} -commitmsg="commit msg unknown" -repository=${PROJECT_NAME} -server_addr=buildrepo:5004 

