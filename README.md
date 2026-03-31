Для установки используйте команды:

git clone https://github.com/Kiseshik/CommentService.git
cd CommentService

Для запуска через докер:
docker build -t comment-service .

Для запуска в in-memory:
docker run -d -p 8080:8080 --name comment-service comment-service

Проверка логов и статуса контейнера:
docker logs comment-service
docker ps

Запуск через докер-композ:
docker-compose up comment-service-memory или docker-compose up comment-service-postgres

Сервер доступен по адресу для режима in-memory: http://localhost:8080
Сервер доступен по адресу для режима postgres: http://localhost:8081

Запуск на фоне:
docker-compose up -d

Быстрый тест на курл:
curl http://localhost:8080