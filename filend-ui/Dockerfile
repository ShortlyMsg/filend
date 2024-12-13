# Filend-UI Dockerfile
FROM node:16-alpine

WORKDIR /app
COPY package.json package-lock.json ./

# Bağımlılıkları yükle
RUN npm install

# Vue uygulamasını derlemek için tüm kaynak kodunu kopyala
COPY . .

# Üretim için Vue uygulamasını derle
RUN npm run build

# Üretim ortamında bir web sunucusu kullan (nginx)
FROM nginx:alpine
COPY --from=0 /app/dist /usr/share/nginx/html

# Nginx'in dinleyeceği portu aç
EXPOSE 9071

# Nginx başlatma komutu
CMD ["nginx", "-g", "daemon off;"]