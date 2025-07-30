#!/bin/bash

CONTAINER_NAME="$1"

if [ -z "$CONTAINER_NAME" ]; then
  echo "Error: No container name provided."
  echo "Usage: ./remove.sh <container-name>"
  exit 1
fi

echo "Removing container: $CONTAINER_NAME"

if docker rm -f "$CONTAINER_NAME" 2>/dev/null; then
  echo "Container '$CONTAINER_NAME' removed successfully."
else
  echo "Container '$CONTAINER_NAME' does not exist or could not be removed."
fi
