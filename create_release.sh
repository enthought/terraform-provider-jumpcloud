#!/bin/bash
if [ -z "$1" ]
    then
    echo "Please provide the release version as the 1st command line arg"
    exit 1
fi

RELEASE_DIR="/tmp/terraform-provider-jumpcloud/releases/$1"

mkdir -p "${RELEASE_DIR}"

echo "Creating source archive"
tar -cvf "${RELEASE_DIR}/source-$1.tar.gz" ./

for GOARCH in 386 amd64
do
  for GOOS in windows darwin linux
  do
    echo "Creating release binary for ${GOOS} / ${GOARCH}"
    # Brute set Go Modules on on chance "auto" does not result in modules.
  	GO111MODULE=on GOOS=${GOOS} GOARCH=${GOARCH} go build -o terraform-provider-jumpcloud
    TARGET_FILENAME="${RELEASE_DIR}/terraform-provider-jumpcloud_${1}_${GOOS}_${GOARCH}.tar.gz"
    tar -cvf "${TARGET_FILENAME}" terraform-provider-jumpcloud
    TARGETS="$TARGETS $TARGET_FILENAME"
  done
done

echo "Creating md5 checksum"
md5sum ${TARGETS} > "${RELEASE_DIR}/checksum-md5.txt"

echo "Creating sh1 checksum"
sha1sum ${TARGETS} > "${RELEASE_DIR}/checksum-sha1.txt"
