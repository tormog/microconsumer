CREATE DATABASE microblogger;

CREATE USER 'example_user'@'%' IDENTIFIED BY 'example_password';
GRANT ALL PRIVILEGES ON microblogger.* TO 'example_user'@'%';
FLUSH PRIVILEGES;

USE microblogger;

CREATE TABLE blogs (
    id VARCHAR(100) NOT NULL,
    source VARCHAR(100) NOT NULL,
    updated DATE,
    data JSON,
    CONSTRAINT blogs_pk PRIMARY KEY (id,source)
);
