FROM mysql:latest

COPY ./sql/schema.sql /docker-entrypoint-initdb.d/01-schema.sql

EXPOSE 3306