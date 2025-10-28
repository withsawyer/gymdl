#!/bin/bash
cd "$(dirname "$0")"
docker rm -f gymdl
docker-compose -f ../docker-compose-local.yml up -d --build