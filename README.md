GoNews Comments Service
=======================
Сервис комментариев для бэкенда GoNews

# Требования

-   golang 1.22

-   docker >=23.0.0

# Запуск

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

# Тесты
TBD
