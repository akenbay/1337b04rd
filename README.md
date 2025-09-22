# 1337b04rd

An anonymous imageboard built with **Go**, **PostgreSQL**, and **S3**.  
Inspired by early internet forums and textboards, `1337b04rd` lets users create threads, post comments, and share images — all without registration.  

---

## Features

- 📝 **Anonymous posting** — no accounts needed; users are tracked via cookies.  
- 💬 **Posts & comments** — create threads, reply to posts or comments.  
- 🖼️ **Image uploads** — stored in an S3-compatible bucket.  
- 👤 **Avatars & nicknames** — automatically assigned via the [Rick & Morty API](https://rickandmortyapi.com/).  
- 🍪 **Session management** — secure cookies with expiration.  
- ⏳ **Post lifecycle** — threads expire after inactivity; archived posts remain view-only.  
- 🏗️ **Hexagonal Architecture** — separation of domain logic and infrastructure.  
- 📜 **Structured logging** — using Go’s `log/slog`.  
- ✅ **Testing** — 20%+ coverage for core functionality.  

---

## Tech Stack

- **Language:** Go (Golang)  
- **Database:** PostgreSQL  
- **Storage:** Own triple-s 
- **Architecture:** Hexagonal (Ports & Adapters)  
- **Other:** net/http, sessions, Rick & Morty API  

---

## Project Structure

- **Domain Layer** — business logic (posts, comments, sessions).  
- **Adapters** — PostgreSQL persistence, S3 image storage, external API clients.  
- **HTTP Layer** — request handlers, middleware, cookie/session management.  

---

## Getting Started

### Prerequisites
- Go 1.22+  
- PostgreSQL  
- Triple-s
- Live server

### Setup

1. Clone the repos:
   ```bash
   git clone https://github.com/akenbay/1337b04rd.git
   git clone https://github.com/akenbay/triple-s.git
   cd 1337b04rd
2. Run the server
   ```bash
   docker-compose up
3. With live server extension open catalog.html in web/templates

### Learning Outcomes
This project helped me practice:
- REST API design
- Cookie-based authentication & sessions
- S3 integration for file storage
- PostgreSQL schema design and SQL transactions
- Hexagonal Architecture for clean separation of concerns
- Logging and error handling in Go
- Writing unit tests with Go’s testing package

### Author 
This project has been created by:

Aibar Kenbay

Contacts:
- email: akenbay@icloud.com
- [GitHub](https://github.com/akenbay/)
- [LinkedIn](https://www.linkedin.com/in/aibar-kenbay-29394b2a4/)
