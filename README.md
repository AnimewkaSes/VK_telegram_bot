# Телеграм бот для сохранения логинов и паролей
## **About**
Данный телеграм бот позволяет сохранять пароли для разных сервисов в одном месте.

Поддерживаются следущие команды:
- /start - выводит приветсвенное сообщение
- /set - добавляет сервис, логин и пароль (сообщения от пользователя удаляются после обработки)
- /get выводит список сервисов, после выбора сервиса выводятся логин и пароль (сообщение с логином и паролем удаляются через 10 секунд)

## **Quick start**
Перед запуском необходимо создать файл env и в .env и поменять значения переменных на свои.

Развертывание приложения происходит в Docker с помощью команды: 
```
docker-compose up -d --build
```
## **Other**
В качестве базы данных используется postgresql.
