FROM mcr.microsoft.com/playwright:v1.41.0-jammy

COPY yarn.lock /yarn.lock
COPY tests /tests

RUN yarn install
