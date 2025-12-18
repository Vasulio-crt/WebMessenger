# Старт
```bash
go run main.go
```

## Просмотр ip
```bash
hostname -i
ip a
```

## Если просит зависимости
```bash
go mod tidy
```


## для себя
```bash
docker volume create mongo_data
docker run -p 27017:27017 -v mongo_data:/data/db mongodb/mongodb-community-server:latest
```