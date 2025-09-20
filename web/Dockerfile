# ---------- build ----------
FROM node:20-alpine AS build
WORKDIR /web
COPY web/package.json web/package-lock.json* ./
RUN npm ci
COPY web ./
RUN npm run build

# ---------- serve ----------
FROM nginx:1.27-alpine
# conf com proxy para a API
COPY web/nginx.conf /etc/nginx/conf.d/default.conf
# build estático
COPY --from=build /web/dist /usr/share/nginx/html
EXPOSE 80