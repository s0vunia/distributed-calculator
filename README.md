# Проект по второму спринту от курса Яндекса Лицея "Разработка на языке Golang"

## _ПО ЛЮБЫМ ВОПРОСАМ ПИШИ МНЕ В ТЕЛЕГРАМ: [@Badimonchik](https://t.me/Badimonchik)_

## Запуск проекта (в корневом каталоге)

1. установить [Make](https://thelinuxcode.com/install-use-make-windows/) (
   опционально), [docker engine](https://docs.docker.com/engine/install/), [docker-compose](https://docs.docker.com/compose/install/)
2. с помощью MakeFile: make build AGENT=3   
   с помощью docker-compose: docker-compose up --scale agent=3 -d --no-recreate --build  
   (вместо трех можно подставить любое число - столько агентов запустится)
3. ждем пару минут (зависит от компьютера и интернет-соединения) пока не запустятся все компоненты системы
## Запросы

* Быстрый импорт запросов с помощью [Postman](docs/Project.postman_collection.json)

1. POST http://localhost:8080/expression - создает expression  
   обязательно должен быть ключ идемпотентности в Header (X-Idempotency-Key) и строка с арифметическим выражением в Body
   form-data (expression)
2. GET http://localhost:8080/expression/{expression_id} - возвращает информацию о выражении
3. GET http://localhost:8080/expressions - возвращает список всех выражений
4. GET http://localhost:8080/agents - возвращает список агентов

## Примеры запросов

1.

``` 
  curl --location 'http://localhost:8080/expression' \
  --header 'X-Idempotency-Key: 1' \
  --form 'expression="2+2*2"'
  ```

```json
a1fd5749-9855-4949-8129-c52bcebbba4f
```

после того, как выражение посчитается

``` 
curl --location 'http://localhost:8080/expression/a1fd5749-9855-4949-8129-c52bcebbba4f'
  ```

```json
{
  "result": 6,
  "id": "a1fd5749-9855-4949-8129-c52bcebbba4f",
  "idempotencyKey": "1",
  "value": "2+2*2",
  "state": "ok"
}
```

2.

``` 
curl --location 'http://localhost:8080/expression' \
--header 'X-Idempotency-Key: 2' \
--form 'expression="2+(2*2)"'
  ```

```json
ad0e375d-0788-4eba-b120-124a4e8b28b7
```

после того, как выражение посчитается

``` 
curl --location 'http://localhost:8080/expression/ad0e375d-0788-4eba-b120-124a4e8b28b7'
  ```

```json
{
  "result": 8,
  "id": "ad0e375d-0788-4eba-b120-124a4e8b28b7",
  "idempotencyKey": "2",
  "value": "(2+2)*2",
  "state": "ok"
}
```

3.

```
curl --location 'http://localhost:8080/expression' \
--header 'X-Idempotency-Key: 3' \
--form 'expression="(2+2)*2))"'
```

```
Invalid expression
```

4. 
```
curl --location 'http://localhost:8080/expression' \
--header 'X-Idempotency-Key: 4' \
--form 'expression="2+2/0"'
```

```json
{
 "result": 0,
 "id": "167935e7-0d86-4f72-991a-079944e3a529",
 "idempotencyKey": "4",
 "value": "2+2/0",
 "state": "error"
}
```

## Структура проекта

## Как работает проект
схема  
описание
## Технологии

## Критерии