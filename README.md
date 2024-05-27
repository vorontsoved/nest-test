# Overview
Приложение для проверка натс

# Шаги для запуска

1.  Вот это
```sh
docker volume create nats1
docker volume create nats2
docker volume create nats3
```

В докер композе должно быть это:

```
version: '3.9'
services:
  nats1:
    image: docker.io/nats:2.9.20
    ports:
      - "4222:4222"
      - "8222:8222"      
    volumes:
      - nats1:/data
    command:
      - "--name=nats1"
      - "--cluster_name=c1"
      - "--cluster=nats://nats1:6222"
      - "--routes=nats-route://nats1:6222,nats-route://nats2:6222,nats-route://nats3:6222"
      - "--http_port=8222"
      - "--js"
      - "--sd=/data"

  nats2:
    image: docker.io/nats:2.9.20
    ports:
      - "4223:4222"
      - "8223:8222"
    volumes:
      - nats2:/data
    command:
      - "--name=nats2"
      - "--cluster_name=c1"
      - "--cluster=nats://nats2:6222"
      - "--routes=nats-route://nats1:6222,nats-route://nats2:6222,nats-route://nats3:6222"
      - "--http_port=8222"
      - "--js"
      - "--sd=/data"

  nats3:
    image: docker.io/nats:2.9.20
    ports:
      - "4224:4222"
      - "8224:8222"
    volumes:
      - nats3:/data
    command:
      - "--name=nats3"
      - "--cluster_name=c1"
      - "--cluster=nats://nats3:6222"
      - "--routes=nats-route://nats1:6222,nats-route://nats2:6222,nats-route://nats3:6222"
      - "--http_port=8222"
      - "--js"
      - "--sd=/data"

volumes:
  nats1:
    external: true
  nats2:
    external: true
  nats3:
    external: true
```

Запускаем кластер:

```
docker-compose up 
```

Проверяем стрим:

```
nats -s localhost:4222 stream ls
No Streams defined
```

Удаляем стрим:

```
nats stream rm
```

2. Переходите в папки публишера и субскрайбера и запускайте мейн файлы
