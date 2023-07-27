# prove

Prove - REST API сервис предоставляет решения алгоритмических задач.


## Запуск сервера.
Обязательно указать DSN базы данных. Рекомендуется PostgreSQL. 
Так же положить папку config/config.json рядом с бинарником.
Если есть сервер удалённый можно запустить docker-compose.yml, 
тем самым поднять DB Postgres и Adminer, в этом случаи adminer будет доступен по :8077

### Подключить DB
При запуске бинарника для флага -d указать ваш DSN
```azure
 prove_<version> -d <DSN>
```

Подключить DB через переменную окружения
```
DATABASE_URI=<DSN> 
```

Проверка работы DB

GET http://localhost:8080/ping
200 Ok
500 Error

## API
* {{port}} = :8080
* {{domain}} = http://localhost{{port}}/api

### Создание группы
* POST {{dmain}}/group 201
Добавить уникальное название группы
```azure
{
    "name":"ИП-200"
}
```

### Все группы
* GET {{domain}}/group 200


### Регистрация пользователя 
<em><small>знаю, что этого не было в тз, но было "как бы ты его писал, работая над энтерпрайз проектом"</small></em>

Создание нового студента, используя его ФИО и т.д.
* POST {{dmain}}/user/register 201
```azure
{
    "email": "bob@mail.ru",
    "login":"bob",
    "password": "mypass",
    "surname": "Штирлиц",
    "name": "Иван",
    "patronymic": "Васильевич",
    "group_id": "d47138a2-2bb9-11ee-9e4d-5b1d4482e44d"
} 
```
### Авторизация 
<em>Кто первый зарегистрировался тот и админ. Админ может добавлять ответы на задачи, по сути он и является продавцом.</em>
* POST {{dmain}}/user/login 200
```azure
{
    "login": "bob",
    "password": "mypass"
} 
```




## Команды Makefile
### Start

Старт prove (скомпилирует и запустит)
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
            
