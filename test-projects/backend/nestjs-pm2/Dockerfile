FROM node:24-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install --no-audit --prefer-offline

COPY . .

RUN npm run build

FROM node:24-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install --only=production --no-audit --prefer-offline

RUN npm install -g pm2

COPY --from=builder /app/dist ./dist

EXPOSE 3000

ENV NODE_ENV=production

CMD ["pm2-runtime", "start", "dist/main.js", "--name", "nestjs-app"]
