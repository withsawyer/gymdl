@echo off
cd /d %~dp0
docker rm -f gymdl
docker-compose -f ../docker-compose-local.yml up -d --build