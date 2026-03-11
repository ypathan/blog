FROM golang:1.26.1-trixie AS builder
WORKDIR /app
COPY . /app
RUN go build -o blog .


FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/blog /app
COPY --from=builder /app/database/init.sql /app
COPY --from=builder /app/static /app/static
EXPOSE 8080 
ENTRYPOINT [ "./blog" ]
