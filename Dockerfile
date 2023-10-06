FROM node:20-alpine AS nodeBuild

WORKDIR /app
COPY vite.config.ts package.json .yarnrc.yml yarn.lock ./
COPY internal/front/ ./internal/front/
COPY .yarn/releases/yarn-3.6.4.cjs ./.yarn/releases/

ENV NPM_CONFIG_PREFIX=/app/.npm-global
# optionally if you want to run npm global bin without specifying path
ENV PATH=$PATH:/app/.npm-global/bin
RUN corepack enable
RUN corepack prepare yarn@stable --activate
RUN yarn set version berry 
RUN yarn 
RUN yarn build

FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && CGO_ENABLED=1 go build -o server ./cmd/api/

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/server .
COPY .env.production ./.env
COPY --from=nodeBuild  /app/internal/front/dist ./web
COPY db/migrations ./db/migrations

CMD [ "/app/server" ]