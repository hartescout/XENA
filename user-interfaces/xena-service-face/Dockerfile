FROM node:16.1.0-alpine

# Create destination directory.
RUN mkdir -p /usr/src/xena-service-face
WORKDIR /usr/src/xena-service-face

# Update and install dependency.
RUN apk update && apk upgrade
RUN apk add git

# Copy the app, note .dockerignore
COPY . /usr/src/xena-service-face/
RUN npm install
RUN npm run build

EXPOSE 3000

ENV NUXT_HOST=0.0.0.0
ENV NUXT_PORT=3000

CMD [ "npm", "start" ]