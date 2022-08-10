# GO MINI PROJECT API CELERATES ACCELERATION PROGRAM

## CAP go miniproject using Echo framework, GORM, and Postgresql

### This Project use 2 domain 
1. Users : Create, Login, Get data
2. Books : Create, Update, Delete, Get Data

### Inside the project
1. Miniproject using ECHO Framework, GORM and Postgresql
2. Using Clean Architecture Pattern Design
3. JWT Token
3. Hasing from bcrypt package
4. Unit Testing

### Users can 
```
1. Create users = localhost:9000/users Method POST 
Bring body format JSON example : 
{
    "name" : "saya",
    "email" : "saya@gmail.com",
    "password" : "saya"
}
response 
{
    "status": 201,
    "messages": "Success create user",
    "data": {
        "id": 35,
        "name": "saya",
        "email": "saya@gmail.com",
        "password": "$2a$10$HaLsYDTwjcPkoag9p4Sjze15MAJPw.OVSQp4QxpvYv9nNNnJwBw3u"
    }
}
2. Loginuser = localhost:9000/auth/login Method POST
bring body format json 
{  
    "email" : "saya@gmail.com",
    "password" : "saya"
}
  
response 
  {
    "status": 200,
    "messages": "Success",
    "data": {
        "name": "saya",
        "email": "saya@gmail.com",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkIjoxNjYwMTQ3MjIxLCJpZCI6IjM1IiwidXNlcm5hbWUiOiJzYXlhQGdtYWlsLmNvbSJ9.2Hqzz_bevI2dsEsBZ_mWVL6x4fuSV7cP08-WIXdCXYY"
    }
}
```

### Books can
```
1. GET books = localhost:9000/books Method GET
bring token from login 
response : 
{
    "status": 200,
    "messages": "Success",
    "data": [
        {
            "id": 1,
            "tittle": "booktest",
            "author": "bookauthor",
            "year": 2002
        }
    ]
}
{
    "status": 200,
    "messages": "Success",
    "data": [
        {
            "id": 2,
            "tittle": "book2",
            "author": "book2",
            "year": 2002
        }
    ]
}

2. GET by ID books = localhost:9000/books/1 Method GET
response
{
    "status": 200,
    "messages": "Success",
    "data": [
        {
            "id": 1,
            "tittle": "booktest",
            "author": "bookauthor",
            "year": 2002
        }
    ]
}

3. UPDATE books = localhost:9000/books/2
response {
    "status": 200,
    "messages": "Success",
    "data": {
        "id": 2,
        "tittle": "updatebook1",
        "author": "updatebook1",
        "year": 2001
    }
}
4. Delete books = localhost:9000/books method delete
response
{
    "status": 200,
    "messages": "Success",
    "data": "Book with id 2 has been deleted"
}
```
