#!/bin/bash

set -e

defaultImage="hello-world"
folder="dumps"

image="${1:-$defaultImage}"
imageFolder="${folder}/${image}"

mkdir -p ${folder}

if [ -d "$imageFolder" ]; then
    echo "Folder '$imageFolder' already exists"
else
    echo "Creating image-specific folder '$imageFolder'..."
    mkdir -p ${imageFolder}
fi

container=$(docker create "$image")

docker export "$container" -o "./${imageFolder}/${image}.tar.gz" > /dev/null

docker inspect -f '{{.Config.Cmd}}' "$image:latest" | tr -d '[]\n' > "${imageFolder}/${image}-cmd"

docker rm "$container" > /dev/null

echo "Image content stored in ${imageFolder}/${image}.tar.gz"
echo "Command configuration stored in ${imageFolder}/${image}-cmd"
echo "Done."