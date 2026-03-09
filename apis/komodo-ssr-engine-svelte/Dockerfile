FROM oven/bun:slim AS build

COPY komodo-forge-sdk-ts /komodo-forge-sdk-ts

WORKDIR /app

COPY komodo-ssr-engine-svelte/package.json komodo-ssr-engine-svelte/bun.lock ./
RUN bun install --frozen-lockfile

COPY komodo-ssr-engine-svelte ./
RUN bun run build

FROM oven/bun:slim AS runtime
WORKDIR /app
COPY --from=build /app/build ./build
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/package.json ./
EXPOSE 7001
CMD ["bun", "build/index.js"]
