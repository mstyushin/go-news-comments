GoNews Comments Service
=======================
Сервис комментариев для бэкенда GoNews

# Конфигурационные параметры

| Параметр         | Описание                           | Значение по-умолчанию                                         |   
|------------------|------------------------------------|---------------------------------------------------------------|
| `http_port`      | порт Comments сервиса              | `8082`                                                        |   
| `db_conn_string` | Строка подключения к СУБД Postgres | `postgres://postgres@localhost:5433/comments?sslmode=disable` |   

# Сборка и запуск

## Требования

-   golang 1.22
-   docker >=23.0.0

Просмотр конфигурации по-умолчанию:

    $ make build
    $ ./bin/go-news-comments -print-config

---

Для быстрого запуска с конфигом по-умолчанию:

    $ make run

Сервер будет запущен на `127.0.0.1:8082`. 

---

Запуск сервера с конфигурацией из файла `config.yaml`:

    $ ./bin/go-news-comments -config config.yaml

---

Остановить сервер:

    $ make clean

---

Показать версию сборки:

    $ ./bin/go-news-comments -version

## Примеры тестовых запросов

    $ curl -XPOST http://127.0.0.1:8082/comments -d '{"author": "Alice", "text": "Hello world2", "article_id": 68}'
    $ curl -XPOST http://127.0.0.1:8082/comments -d '{"author": "Bob", "Text": "Answering to your hello", "article_id": 68, "parent_id": 1}'

# Тесты
Полный прогон имеющихся тестов:

    $ make test

Для удаления тестовой базы:

    $ make clean
