#!/usr/bin/env bash

set -e

cd web
npm install
npm run-script build
cp -r build/* ../orgnetsim/web/

cd ../orgnetsim
go build