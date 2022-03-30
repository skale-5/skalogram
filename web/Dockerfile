FROM golang:1.18-bullseye AS build
RUN mkdir -p /skalogram-web
WORKDIR /skalogram-web
COPY ./ ./
RUN go build -o out ./cmd/server

FROM gcr.io/distroless/base-debian11
COPY --from=build /skalogram-web/out /usr/local/bin/skalogram-web
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/skalogram-web"]