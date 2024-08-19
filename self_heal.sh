#!/bin/bash

if ! command -v nvidia-smi &> /dev/null; then
    echo "$(date): not found gpu"
    exit 0
fi

while true; do
    if ! nvidia-smi > /dev/null 2>&1; then
        echo "$(date): nvidia-smi not use, retrying..."
        kill -9 resource-exporter
    fi
    sleep 60
done
