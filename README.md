
# Ресайзер картинок

Сервис представляет собой web-сервер (прокси), загружающий изображения,
масштабирующий/обрезающий их до нужного формата и возвращающий пользователю.

Передаем URL адрес картинки и новые размеры, получаем картинку измененного размера. Переданные http заголовки проксируются и передаются к целевому серверу.
Картинка кешируется в minio

## API Reference

```http
  GET /${width}/${height}/${url}
```

| Параметр | Тип     | Описание                |
| :-------- | :------- | :------------------------- |
| `width`   | `int`    | **Required**. Ширина до которой изменяем картинку |
| `height`  | `int`    | **Required**. Длина до которой изменяем картинку |
| `url`     | `string` | **Required**. URL откуда качаем картинку, без "http://"|


#### Например:
```http
http://localhost:8082/150/100/habrastorage.org/r/w1560/getpro/habr/post_images/167/521/f9b/167521f9b392c45594e33f659165bdbb.png
```

## Deployment



```bash
  docker-compose up
```


## Environment Variables

Переменные окружения, которые есть в проекте

```
ENV=local # (local, dev, prod)
MINIO_ENDPOINT=minio:9000
MINIO_PORT=9000
MINIO_ACCESSKEY=minioadmin
MINIO_SECRETKEY=minioadmin
MINIO_BUCKET=dev-minio
HTTP_SERVER_ADDRESS=localhost:8082
HTTP_SERVER_TIMEOUT=4s
HTTP_SERVER_IDLE_TIMEOUT=60s
```

