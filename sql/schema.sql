# DROP DATABASE IF EXISTS userdb;

# CREATE DATABASE userdb;

USE userdb;

CREATE TABLE users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    age INTEGER NOT NULL
)