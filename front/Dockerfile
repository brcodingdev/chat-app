FROM node:20.3.1

WORKDIR /
COPY ./front /
COPY ./front/package.json ./

WORKDIR /app

RUN npm install

COPY . .

RUN npm run build
EXPOSE 3000

CMD ["npm", "start"]