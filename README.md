## gin-rest-api-sample
Este proyecto-base de API REST esta escrito en GOLANG y usa como base de datos MariaDB usando como orm GORM

Este proyecto cuenta con los siguientes paquetes:

- REST API servidor con [Gin Framework](https://github.com/gin-gonic/gin)
- Database integrado usando [GORM](http://gorm.io/)
- Live Reload usando [codegangsta/gin](https://github.com/codegangsta/gin)
- JWT Token based Authentication
- [Supported REST API Documentation](https://documenter.getpostman.com/view/723994/RWTeVNA4) (Postman)


## Iniciar el Proyecto

```
$ go run main.go
```

Ahora compilaremos antes de subir al servidor:

```
$ go build main.go
$ ./main
```

Usamos live-reloading en entorno de desarrollo

```
$ ./scripts/start-dev
```
