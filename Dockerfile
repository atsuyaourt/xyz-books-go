FROM node:iron-alpine3.19 AS node-pnpm
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app

COPY . .
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm exec postcss /app/internal/views/style.css -o /app/internal/assets/style.css

FROM golang:1.22-alpine3.19 AS go
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=go /app/server .
COPY .env.production ./.env
COPY --from=node-pnpm  /app/internal/assets ./internal/assets/
COPY internal/db/migrations ./db/migrations

CMD [ "/app/server" ]
