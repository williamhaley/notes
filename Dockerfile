FROM golang:1.23.4-bookworm AS builder

FROM builder AS web-app-builder

WORKDIR /web-app

RUN curl -fsSL https://deb.nodesource.com/setup_23.x -o nodesource_setup.sh
RUN bash nodesource_setup.sh
RUN apt-get install -y nodejs

COPY ./ /web-app/

RUN corepack enable
RUN corepack install
RUN pnpm install
RUN pnpm build

FROM builder AS run-time

WORKDIR /usr/src/app

COPY ./api/go.mod ./api/go.sum ./
RUN go mod download && go mod verify

COPY ./api .

COPY --from=web-app-builder /web-app/dist/ /usr/src/app/www/public_html/

RUN go build -v -o /usr/local/bin/app ./cmd/...

CMD ["app"]
