# --- Build stage ---
FROM node:20-bookworm AS builder

WORKDIR /app

COPY package.json package-lock.json ./

# Clean install with optional dependencies included
RUN npm ci --include=optional

COPY . .

ARG VITE_API_BASE_URL=http://localhost:8080
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

RUN npm run build

# --- Runtime stage: serve the built static files with nginx ---
FROM nginx:1.27-alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
    CMD wget -q -O - http://localhost:80 || exit 1

CMD ["nginx", "-g", "daemon off;"]
