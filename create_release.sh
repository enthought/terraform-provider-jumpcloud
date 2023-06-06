#!/bin/bash
if [ -z "$1" ]
    then
    echo "Please provide the release version as the 1st command line arg"
    exit 1
fi

RELEASE_DIR="/tmp/terraform-provider-jumpcloud/releases/$1"

echo "Creating source archive"
tar -cvf "${RELEASE_DIR}/source-$1.tar.gz" ./

mkdir -p "${RELEASE_DIR}"
for goarch in amd64 arm64
do
  for goos in darwin
  do
    echo "creating release binary for ${goos} / ${goarch}"
  	env goos=${goos} goarch=${goarch} go build -o terraform-provider-jumpcloud
    target_filename="${release_dir}/terraform-provider-jumpcloud_${1}_${goos}_${goarch}.tar.gz"
    tar -cvf "${target_filename}" terraform-provider-jumpcloud
    targets="$targets $target_filename"
  done
done

for goarch in amd64
do
  for goos in linux
  do
    echo "creating release binary for ${goos} / ${goarch}"
  	env goos=${goos} goarch=${goarch} go build -o terraform-provider-jumpcloud
    target_filename="${release_dir}/terraform-provider-jumpcloud_${1}_${goos}_${goarch}.tar.gz"
    tar -cvf "${target_filename}" terraform-provider-jumpcloud
    targets="$targets $target_filename"
  done
done




echo "Creating md5 checksum"
md5sum ${TARGETS} > "${RELEASE_DIR}/checksum-md5.txt"

echo "Creating sh1 checksum"
sha1sum ${TARGETS} > "${RELEASE_DIR}/checksum-sha1.txt"
