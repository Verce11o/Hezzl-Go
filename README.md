# Hezzl-Go
Тестовое задание от Hezzl.com
## О проекте
- Проект был создан в чистой архитектуре. 
- Слои общаются между собой через интерфейсы.
- В NATS был использован JetStream для гарантированной доставки сообщений.
- Есть поддержка линтера golangci-lint
- Есть документация swagger
## Инструкция к запуску
Для запуска проекта требуется docker-compose.

`docker-compose up -d`

Далее запустите миграции. 

`make postgres-up`
`make clickhouse-up`

## API
Endpoint = `http://localhost:3000/api/v1`
