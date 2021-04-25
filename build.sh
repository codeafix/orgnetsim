#!/usr/bin/env bash

set -e

cd web
npm run-script build
cp -r build/* ../orgnetsim/web/

cd ../orgnetsim
go build