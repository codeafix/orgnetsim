#!/usr/bin/env bash

set -e

cd web
npm install
npm run build
cp -r dist/* ../orgnetsim/web/

cd ../orgnetsim
go build