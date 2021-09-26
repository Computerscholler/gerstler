#!/bin/bash
docker-compose up -d
cd backend/gerstler
go run cmd/main.go


