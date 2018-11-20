FROM node:10-alpine

WORKDIR /ng-app/

COPY ./package.json ./package-lock.json /ng-app/
RUN npm install
COPY . .
RUN npm run build
RUN ls -la /ng-app/dist

