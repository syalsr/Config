# Distributed config
## Описание
Реализовать сервис, который позволит динамически управлять конфигурацией приложений. Доступ к сервису должен осуществляться с помощью API. Конфигурация может храниться в любом источнике данных, будь то файл на диске, либо база данных. Для удобства интеграции с сервисом может быть реализована клиентская библиотека.

## Реализация
Для выполнения задания использовал gRPC, grpc gateway, Postgresql, docker-compose

## Запуск
```
make run
```

## Проблемы
С изначальной структурой конфига были проблемы
```
{
    "service": "managed-k8s",
    "data": [
        {"key1": "value1"},
        {"key2": "value2"}
    ]
}
```
В протобафе я пытался сделать такую структуру данных:
```
message RequestConfig{
  string service = 1;
  repeated Data data = 2;
}
message Data{
  map<string, string> key = 1;
}
```

Но данные не приходили, поэтому на выбор было 2 варианта - костыльный
```
message Request{
  string service = 1;
  repeated Data data = 2;
}
message Data{
  string key1 = 1;
  string key2 = 2;
}
```
Тогда любой конфиг должен иметь только 2 ключа/значения. Либо изменить саму структуру конфига, чтобы я мог изменять количество ключей/значений в конфигах:
```
{
    "service": "test1",
    "data": {
        "key1": "value1",
        "key2": "value2",
        "key3": "dsfsd"
    }
}
```
```
message RequestConfig{
  string service = 1;
  map<string, string> data = 2;
}
```
