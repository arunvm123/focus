FROM golang:alpine AS build-env

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# All these steps will be cached
RUN mkdir /travail-backend
WORKDIR /travail-backend
# <- COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN go build -o server .
# <- Second step to build minimal image
FROM scratch
COPY --from=build-env /travail-backend/certs /certs
COPY --from=build-env /travail-backend/push_notification/fcm/travail-7f7b9-firebase-adminsdk-v5arf-c7dd3d30d3.json /
COPY --from=build-env /travail-backend/server /
# COPY --from=build-env /travail-backend/config.yaml /
EXPOSE 5000
ENTRYPOINT ["./server","-config-env","true"]