#!/bin/bash

DOCKERFILE_PATH="./images/m1-number-printer/Dockerfile"

BUILD_CONTEXT="./images/m1-number-printer/"

IMAGE_NAME="m1-number-printer-image"
IMAGE_TAG="1.0"

echo "Building Docker image $IMAGE_NAME:$IMAGE_TAG from $DOCKERFILE_PATH in $BUILD_CONTEXT..."
docker build -f "$DOCKERFILE_PATH" -t "$IMAGE_NAME:$IMAGE_TAG" "$BUILD_CONTEXT"

if docker build -f "$DOCKERFILE_PATH" -t "$IMAGE_NAME:$IMAGE_TAG" "$BUILD_CONTEXT"; then
    echo "Docker image $IMAGE_NAME:$IMAGE_TAG was built successfully."
else
    echo "Failed to build Docker image."
fi