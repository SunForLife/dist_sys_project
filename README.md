## Distributed online store

Проект представляет собой онлайн магазин состоящий из нескольких микросервисов написанных на Go.

Этими микросервисами являются: основной бэкэнд обрабатывающий запросы с логикой для товаров, сервис авторизации поддерживающий работу через access и refresh токены, сервис подтверждения регистрации и операций черес смс, и сервис загрузки в базу данных продуктов через готовый csv дамп посредствам потоковой передачи данных.

Postman можно найти в папке api.

### Диаграмма сервиса

![diagram](./diagram.jpg)

### Запуск сервера

`sudo docker-compose up --build`

