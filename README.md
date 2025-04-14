# TaskManager API 🗂️

A clean and scalable task management backend API built in **Go** using the **Gin** web framework, following professional software architecture principles (MVC + Clean Architecture).

## 🔧 Features

- JWT-based Authentication (Register/Login)
- Secure Password Hashing with bcrypt
- CRUD operations for Users and Tasks
- Layered architecture (Controllers, Services, Repositories)
- PostgreSQL integration using GORM
- Middleware for Authorization
- Structured project bootstrap with `cmd/server`

## 📁 Project Structure

TaskManager/ ├── cmd/ │ └── server/ # Main entrypoint ├── internal/ │ ├── controllers/ │ ├── services/ │ ├── repositories/ │ ├── models/ │ ├── routes/ │ ├── config/ │ └── middleware/ ├── pkg/ │ └── utils/ # Password utils, JWT, etc. ├── bootstrap/ # App initialization ├── go.mod ├── .env

## 📦 Technologies Used

- **Go (Golang)**
- **Gin** - Web Framework
- **GORM** - ORM for PostgreSQL
- **JWT** - JSON Web Token Auth
- **bcrypt** - Secure password hashing

## 🚀 Getting Started

1. Clone the repository:
   ```bash
   git clone github.com/bishalcode860/TaskManagerProjectGo
   cd TaskManagerGo
