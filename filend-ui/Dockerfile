# Filend-UI Dockerfile
FROM node:20

WORKDIR /app
COPY package.json package-lock.json ./

# Bağımlılıkları yükle
RUN npm install

# Vue uygulamasını derlemek için tüm kaynak kodunu kopyala
COPY . .

# Üretim için Vue uygulamasını derle
#RUN npm run build

# Üretim ortamında bir web sunucusu kullan (nginx)
# FROM nginx:alpine
# COPY --from=0 /app/dist /usr/share/nginx/html
# CMD ["nginx", "-g", "daemon off;"]

EXPOSE 9071

CMD ["npm", "run", "serve"]