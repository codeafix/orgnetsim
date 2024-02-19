#!/usr/bin/env bash

set -e

cd web
npm install
tsc
vite build
cp -r build/* ../orgnetsim/web/

cd ../orgnetsim
go build