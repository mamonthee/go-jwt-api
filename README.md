# Articles and Authors API

This project implemented a RESTful API built using the Gin framework in Go and SQLite3 as the database. It includes user authentication with JWT tokens, author management, and article CRUD operations.
---

## API Endpoints

#### Register a User
```bash
curl -X POST http://localhost:9000/register \
     -d '{"user_name":"mary","email":"mary@gmail.com", "password":"mary123456"}'
```

#### Login
```bash
curl -X POST http://localhost:9000/login \
     -d '{"email":"john@gmail.com","password":"john123456"}'
```

---

### **Authors**

#### Update Author
```bash
curl -X PUT http://localhost:9000/author/update \
     -H "Authorization: Bearer <YOUR_TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"user_name":"John","email":"john@gmail.com"}'
```

#### Deactivate Author
```bash
curl -X PUT http://localhost:9000/author/deactivate \
     -H "Authorization: Bearer <YOUR_TOKEN>" \
     -H "Content-Type: application/json"
```

---

### **Articles**

#### Get Articles
```bash
curl -X GET http://localhost:9000/articles/ \
     -H "Authorization: Bearer <YOUR_TOKEN>"
```

#### Create Article
```bash
curl -X POST http://localhost:9000/articles/ \
     -H "Authorization: Bearer <YOUR_TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"title":"Third Article Of John","description":"This is an example article."}'
```

#### Update Article
```bash
curl -X PUT http://localhost:9000/articles/5 \
     -H "Authorization: Bearer <YOUR_TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"title":"Updated Article Title","description":"Updated article description."}'
```

#### Delete Article
```bash
curl -X DELETE http://localhost:9000/articles/5 \
     -H "Authorization: Bearer <YOUR_TOKEN>"
```