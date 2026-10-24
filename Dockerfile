FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=$VERSION" -o /ledger .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /ledger /ledger
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s CMD ["/ledger", "healthcheck"]
ENTRYPOINT ["/ledger"]
