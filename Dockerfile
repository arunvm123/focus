FROM golang:1.12.7 as build-env
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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/travail-backend
# <- Second step to build minimal image
FROM scratch
COPY --from=build-env /travail-backend/config.yaml .
COPY --from=build-env /go/bin/travail-backend .
EXPOSE 5000
ENTRYPOINT ["./travail-backend"]