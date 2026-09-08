# Analyst Hall: Vue dist + Go API on 8080. Secrets stay in runtime env, not in the image.
FROM node:22-alpine AS web
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web ./
RUN npm run build \
  && test -f dist/media/landing/hero-desk.jpg \
  && test -f dist/media/landing/bank-pay.jpg

FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY --from=web /web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o /academy ./cmd/academy

FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates wget \
  && rm -rf /var/lib/apt/lists/*
COPY --from=build /academy /app/academy
COPY --from=build /src/web/dist /app/web/dist
ENV PORT=8080
ENV HOST=0.0.0.0
ENV WEB_DIR=/app/web/dist
ENV DATA_DIR=/app/data
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=12s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/api/health >/dev/null || exit 1
CMD ["/app/academy"]
