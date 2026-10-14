# Production image — build context is this folder (fe/), so fe/ can be built and
# deployed on its own. NEXT_PUBLIC_* values are inlined at build time.
FROM node:22-alpine AS builder

WORKDIR /app

RUN corepack enable && corepack prepare pnpm@10.12.1 --activate

ARG NEXT_PUBLIC_API_BASE_URL=http://localhost/api
ARG NEXT_PUBLIC_APP_NAME=System App
ARG NEXT_PUBLIC_USE_MOCK_API=false
ARG NEXT_PUBLIC_SITE_URL=http://localhost

ENV NEXT_PUBLIC_API_BASE_URL=$NEXT_PUBLIC_API_BASE_URL
ENV NEXT_PUBLIC_APP_NAME=$NEXT_PUBLIC_APP_NAME
ENV NEXT_PUBLIC_USE_MOCK_API=$NEXT_PUBLIC_USE_MOCK_API
ENV NEXT_PUBLIC_SITE_URL=$NEXT_PUBLIC_SITE_URL

COPY package.json pnpm-lock.yaml .npmrc ./
RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm build

FROM node:22-alpine AS runner

WORKDIR /app

RUN corepack enable && corepack prepare pnpm@10.12.1 --activate

ENV NODE_ENV=production

COPY package.json pnpm-lock.yaml .npmrc ./
RUN pnpm install --prod --frozen-lockfile

COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY next.config.ts postcss.config.mjs tsconfig.json ./
COPY src ./src

EXPOSE 3000

CMD ["pnpm", "start", "--hostname", "0.0.0.0", "--port", "3000"]
