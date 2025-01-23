# Hotel reservation backenmd

## Project outline
- Users -> book room from a hotel
- Admins -> going to check reservations/bookings
- Authentication and authorization -> JWT tokens
- Hotels -> CRUD API -> JSON
- Rooms -> CRUD API -> JSON
- Scripts -> database managment -> seeding, migration

## Resources 
### Mongodb driver
Documentation
```
https://www.mongodb.com/docs/drivers/go/current/quick-start/
```

Installing mongodb client
```
go get go.mongodb.org/mongo-driver/mongo
```

### gofiber 
Documentation
```
https://gofiber.io
```

Installing gofiber
```
go get github.com/gofiber/fiber/v2
```

## Docker
### Installing mongodb as a Docker container
```
docker run --name mongodb -d mongo:latest -p 27017:27017
```