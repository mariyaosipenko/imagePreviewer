## Compile stage
#FROM golang:1.22 AS build-env
#
## Work directory
#WORKDIR /image-resize
#
## Installing dependencies
#COPY go.mod go.sum ./
#RUN go mod download
#
## Copying all the files
#COPY . .
#
## Starting our application
#CMD ["go", "run", "image-resize/cmd/img-prev"]
#
## Exposing server port
#EXPOSE 8082


# Latest golang image on apline linux
FROM golang:1.22 as builder

## Install libvips and necessary dependencies
#RUN apk update \
#    && apk upgrade \
#    && apk add \
#    build-base \
#    vips-dev \
#    && rm -rf /var/cache/apk/*

#RUN apt-get update -qq && apt -y upgrade && \
#    apt-get install -y libvips42

# Work directory
WORKDIR /image-resize

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copying all the files
COPY . .

# Building the application
RUN go  build -o image-resize image-resize/cmd

#CMD ["sqlite3", "/storage/storage.db"]

## Fetching the latest nginx image
#FROM alpine:3.16 as production
#
## Certificates
#RUN apk add --no-cache ca-certificates
#
## Copying built assets from builder
#COPY --from=builder image-resize .
#
# Starting our application
CMD ./image-resize

# Exposing server port
EXPOSE 8082