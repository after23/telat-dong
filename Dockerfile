FROM golang:1.20.4-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
RUN touch app.env
RUN echo BOT_TOKEN=SANGATRAHASIA >> app.env
RUN echo PLAYGROUND_ID=123 >> app.env


CMD [ "/app/main" ]