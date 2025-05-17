FROM node:18-alpine

WORKDIR /app

COPY package*.json ./

RUN rm -rf node_modules package-lock.json && \
    npm cache clean --force && \
    npm install && \
    npm install esbuild@0.21.5 --save-dev

COPY . .

RUN npm run build

FROM nginx:alpine

COPY --from=0 /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]