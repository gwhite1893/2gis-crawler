# 2gis-crawler

Сервис предоставляет REST API для для опроса ресурсов.
В данной редакции для каждого урла получает контент ответа и возвращает 
значение тега <title>.

## Установка

```
    go get github.com/gwhite1893/2gis-crawler
```

## Swagger

```
    http://{host}:{port}/api/crawler/v1/swagger/
```

## Использование

- В файле config.yaml указать требуемые параметры опроса и веб-сервера
- Запустить проект 
- Сформировать необходимый список ссылок на ресурсы и передать его в хендлер, например:

``` 
    сurl --location --request POST 'http://{host}:{port}/api/crawler/v1/sources/poll' \
  --header 'Content-Type: application/json' \
  --data '{"data":["http:\\mts1.ru","http://ya.ru", "http://google.com","https://mts.ru", "http://httpstat.us/200?sleep=5000"]}'
```

