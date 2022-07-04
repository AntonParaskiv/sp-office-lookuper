# Accept the Go version for the image to be set as a build argument.
FROM golang:1.18.1-alpine3.15 AS build

RUN apk --no-cache add openssl-dev curl git openssh \
    gcc musl-dev linux-headers util-linux make

## Precompile the entire go standard library into the first Docker cache layer: useful for other projects too!
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go install -v -installsuffix cgo -a std

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user \
    && echo 'user:x:65534:65534::/:' > /user/passwd \
    && echo 'user:x:65534:' > /user/group

ENV GO111MODULE on

WORKDIR /src

COPY go.* ./

# Get and precompile third party libraries,
# See issues https://github.com/golang/go/issues/27719.
## Reusing previous go build cache by `--mount`.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod graph | awk '$1 !~ /@/ { print $2 }' | xargs -r go get -x && \
    go list -m -f '{{ if not .Main }}{{ .Path }}/...@{{ .Version }}{{ end }}' all | tail -n +2 | \
    CGO_ENABLED=1 GOOS=linux xargs go build -v -installsuffix cgo -i; echo done

## Lint and test service code.
COPY . .

# Reusing the linter cache and go cache for run linter and go tests.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    go test -v ./...

# Compile! Should only compile our sources since everything else is precompiled.
ARG RACE
ARG CGO
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir bin && \
    GOOS=linux CGO_ENABLED=${CGO} go build ${RACE} -v -installsuffix cgo -o ./bin \
    -ldflags "-linkmode external -extldflags -static -s -w" \
    ./cmd/api

FROM alpine:3.15

ENV APP_ROOT /opt/user
RUN mkdir -p $APP_ROOT
ENV PATH "$PATH:$APP_ROOT"
WORKDIR $APP_ROOT

RUN adduser -D -u 1111 user
EXPOSE 8080
EXPOSE 8090
COPY --from=build /src/bin $APP_ROOT/
USER user
CMD ["api"]
