FROM node:10.21.0-alpine3.11 as react-build
ARG CI=true
WORKDIR /app
COPY . ./
RUN yarn
RUN yarn test
RUN yarn build

FROM nginx:1.17
COPY --from=react-build /app/build /usr/share/nginx/html

