# Тестовое задание для OZON

Для реализации graphql использована библиотека gqlgen(см. директорию graph)

Добавлена авторизация с применением jwt-токенов.

Структура пользователя:
```Go
type User struct{
	username string `json:"username"`
	password string `json:"password"`
	isAdmin bool `json:"isAdmin,omitempty"`
}
```

Методы: 
- /login 
    Запрос:
    ```JSON
    {
      "username":"username",
      "password":"password"
    }
    ```
  
- /logout
- /register
  Запрос:
    ```JSON
    {
      "username":"username",
      "password":"password"
    }
    ```
Запросы на публикацию и комментирование недоступны для неавторизованных пользователей
Запросы на получение постов и комментариев не требует авторизации

Админ-пользователь: "admin", "admin"

Схема таблицы для Postgres расположена в директории db.

Схема graphql:
```GraphQL
type Post{
    id: ID!
    author: String!
    name: String!
    description: String!
    content: String!
    comments_allowed: Boolean!
    createdAt: Time!
    updatedAt: Time!
}

type Comment{
    id: ID!
    post_id: ID!
    # Replied comment id
    replies_to: ID!
    author: String!
    text: String!
    createdAt: Time!
    updatedAt: Time!
}

type User{
    username: String!
    password: String!
    isAdmin: Boolean!
}

input NewPost{
    name: String!
    description: String!
    content: String!
    comments_allowed: Boolean!
}

input NewComment{
    post_id: ID!
    replies_to: ID!
    text: String!
}

input NewUser{
    username: String!
    password: String!
}

type Mutation{
    createPost(input: NewPost): ID!
    createComment(input: NewComment): ID!
}

type Query{
    Posts(limit: Int = 25, offset: Int = 0): [Post!]!
    Comments(post_id: ID): [Comment!]!
}

scalar Time
```