FROM oven/bun:1-alpine AS base
WORKDIR /app

FROM base AS dev
COPY package.json bun.lock ./
RUN bun install
COPY . .
CMD ["bun", "run", "dev", "--host"]

FROM base AS build
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile
COPY . .
RUN bun run build

FROM base AS prod
COPY --from=build /app/build ./build
COPY --from=build /app/package.json ./
COPY --from=build /app/bun.lock ./
RUN bun install --production --frozen-lockfile
ENV NODE_ENV=production
CMD ["bun", "run", "build/index.js"]