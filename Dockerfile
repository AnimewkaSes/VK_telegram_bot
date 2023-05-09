# Сборка приложения
FROM golang:latest as builder

WORKDIR ./app

COPY . .
RUN echo * && go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /main main.go

# Копируем собранное приложение и добавляем сертификат
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]