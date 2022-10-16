# Maria DB setup
Based on https://github.com/lindycoder/prepopulated-mysql-container-example

1. Build

    docker build -t mariadb-test .

2. Run

    docker run -d --rm -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password --name mariadb mariadb-test
