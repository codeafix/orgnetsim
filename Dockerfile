FROM node:latest as webapp
WORKDIR /app
COPY ./web .
RUN npm install
RUN npm run build

FROM golang:latest as appapi
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -C orgnetsim -tags netgo -o /usr/local/bin/app ./...

FROM alpine:latest
ARG UID
ARG GID

RUN mkdir -p /app/dist
RUN mkdir -p /app/bin
COPY --from=webapp /app/dist /app/dist
COPY --from=appapi /usr/local/bin/app /app/bin/orgnetsim

RUN getent group $GID || addgroup -g $GID defgrp && \
    adduser -u $UID -g $GID --disabled-password --gecos "" default
USER default
RUN mkdir -p /tmp/data
COPY --chown=default:default <<EOF /tmp/data/sims.json
{"simulations":[],"notes":"A simulator for Organisational Networks. The simulator is created from a Network of Agents. The Network itself can be any arbitrary graph and contains a collection of Agents and a collection of links between those Agents. The simulator uses Colors to represent competing ideas on the Network. The default Color for an Agent is Grey. During a simulation Agents interact and decide whether or not to update their Color."}
EOF

ENTRYPOINT [ "/app/bin/orgnetsim" ]
CMD ["serve", "/tmp/data", "-s", "/app/dist"]
EXPOSE 8080