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

После запуска сервера база данных, девственно чистая, поэтому начинать нужно с создания группы.

### Создание группы

* POST {{dmain}}/group 201
Добавить уникальное название группы
```azure
{
    "name":"ИП-200"
}
```

### Все группы
Здесь можно найти id группы
* GET {{domain}}/group 200


### Регистрация пользователя 
<em><small>знаю, что этого не было в тз, но было "как бы ты его писал, работая над энтерпрайз проектом"</small></em>

Создание нового студента, используя его ФИО и т.д. 

<em>Обязательно указать group_id (можно посмотреть получив список всех групп).</em>
* POST {{domain}}/user/register 201
```azure
{
    "email": "bob@mail.ru",
    "login":"bob",
    "password": "mypass",
    "surname": "Штирлиц",
    "name": "Иван",
    "patronymic": "Васильевич",
    "group_id": "<group_id>"
} 
```
### Авторизация 
<em>Кто первый зарегистрировался тот и админ. Админ может добавлять ответы на задачи, по сути он и является продавцом.</em>
* POST {{domain}}/user/login 200
```azure
{
    "login": "bob",
    "password": "mypass"
} 
```

### Текущий кредит
Возможно получить ответ на задачу пока ваш кредит не достиг 1000 
* GET {{domain}}/user/balance

### Пополнить баланс
* POST {{domain}}/user/balance
```azure
{
    "current":100.30
}
```

### Создать задачу
* POST {{domain}}/task
```azure
{
    "name":"Найти все пропущенные числа.", 
    "description":"Дан неотсортированный массив из N чисел от 1 до N,при этом несколько чисел из диапазона [1, N] пропущено, а некоторые присутствуют дважды.",
    "price":"100.87"
}
```

### Получить список задач
* GET {{domain}}/task
response:
```azure
[
    {
        "task_id": "<task_id>",
        "name": "Найти все пропущенные числа.",
        "description": "Дан неотсортированный массив из N чисел от 1 до N,при этом несколько чисел из диапазона [1, N] пропущено, а некоторые присутствуют дважды.",
        "price": "100.87",
        "created": "2023-07-27T16:41:19+03:00"
    }
]
```
### Получить конкретную задачу
* GET {{domain}}/task/{{task_id}}
response:
```azure
{
    "task_id": "<task_id>",
    "name": "Найти все пропущенные числа.",
    "description": "Дан неотсортированный массив из N чисел от 1 до N,при этом несколько чисел из диапазона [1, N] пропущено, а некоторые присутствуют дважды.",
    "price": "100.87",
    "created": "2023-07-27T16:41:19+03:00"
}
```
### Получить ответ по задачи
* POST {{domain}}/task/solution

```azure
{
    "task_id":"<task_id>",
    "data":[1, 1, 2, 3, 5, 5, 7, 8, 9, 9]
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
            
