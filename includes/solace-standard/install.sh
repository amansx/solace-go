#!/usr/bin/env bash

cd includes/solace-standard
cat x* > solace-pubsub-standard-9.5.0.30-docker.tar.gz
docker load -i *.tar.gz
rm solace-pubsub-standard-9.5.0.30-docker.tar.gz
