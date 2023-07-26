# prove

Prove - REST API сервис предоставляет решения алгоритмических задач.


## Запуск сервера из бинарника.
Обязательно указать DSN базы данных. Рекомендуется PostgreSQL. 
Так же положить папку config/config.json рядом с бинарником.
Если есть сервер удалённый можно запустить docker-compose.yml, 
тем самым поднять DB Postgres и Adminer, в этом случаи adminer будет доступен по :8077
* Можно подключиться к demo серверу сейчас так:
* prove -d postgres://bob:mypasswd@194.67.104.167:5432/prove?sslmode=disable

### Linux
Подключение DB через флаг -d
```
./bin/prove_<здесь версия>_linux_amd6 -d postgres://prove:DBprove-2023@127.0.0.1:5432/prove?sslmode=disable
```

Возможно так же подключить DB через переменную окружения
```
DATABASE_URI=postgres://prove:DBprove-2023@127.0.0.1:5432/prove?sslmode=disable ./bin/prove_<здесь версия>_linux_amd6 
```




## Команды Makefile
### Start

Старт prove
```azure
make build_d
```

Сборка для трех os
```azure
make build_a
```

HELP по флагам
```azure
make short_h
```
            
