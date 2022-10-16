# Maria DB setup
Based on https://github.com/lindycoder/prepopulated-mysql-container-example

## Build

    docker build -t mariadb-test .

## Run

    docker run -d --rm -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password --name mariadb mariadb-test
