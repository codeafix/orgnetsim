#!/usr/bin/env bash

set -e

cd web
npm install
vite build
cp -r dist/* ../orgnetsim/web/

cd ../orgnetsim
go build