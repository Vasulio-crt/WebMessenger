# Старт
```bash
go run main.go
```

## Для себя
```bash
docker volume create mongo_data
docker run -p 27017:27017 -v mongo_data:/data/db mongodb/mongodb-community-server:latest
```

```bash
hostname -i
ip a
```

```bash
go mod tidy
```