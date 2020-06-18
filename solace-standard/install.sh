#!/usr/bin/env bash

cat x* > solace-pubsub-standard-9.5.0.30-docker.tar.gz
docker load -i solace-standard/*.tar.gz
