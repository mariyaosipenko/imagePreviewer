# Latest golang image on apline linux
FROM golang:1.22 as builder

# Work directory
WORKDIR /image-resize

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copying all the files
COPY . .

# Building the application
RUN go  build -o image-resize image-resize/cmd

CMD ./image-resize

# Exposing server port
EXPOSE 8082