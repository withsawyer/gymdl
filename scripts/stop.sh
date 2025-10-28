#!/bin/bash
cd "$(dirname "$0")"
docker-compose -f ../docker-compose-local.yml down