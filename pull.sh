#!/bin/bash

set -e

defaultImage="hello-world"
folder="dumps"

image="${1:-$defaultImage}"

echo "Creating a folder '$folder' for storing files..."
mkdir -p ${folder} || handle_error "Failed to create ${folder} directory"

echo "Creating temporary container from image '$image'..."
container=$(docker create "$image")

echo "Extracting image '$image' to '${folder}/${image}.tar.gz'..."
docker export "$container" -o "./${folder}/${image}.tar.gz" > /dev/null

echo "Extracting default command configuration..."
docker inspect -f '{{.Config.Cmd}}' "$image:latest" | tr -d '[]\n' > "./${folder}/${image}-cmd"

echo "Cleaning up temporary container..."
docker rm "$container" > /dev/null

echo "Image content stored in ${folder}/${image}.tar.gz"
echo "Command configuration stored in ${folder}/${image}-cmd"
echo "Done."