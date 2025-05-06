FROM node:20.19.1
LABEL authors="cilikube"

WORKDIR /app
COPY package.json ./

RUN npm install

COPY . .

EXPOSE 8888
ENTRYPOINT ["npm", "run", "dev"]
