# 1337b04rd  
_Anonymous imageboard for hackers — built with Go, PostgreSQL, and S3._

## 📖 Overview

`1337b04rd` is a minimalistic yet functional imageboard (think old-school forums) with:
- Anonymous posting (no registration)
- Image uploads via S3
- Unique avatars from The Rick and Morty API
- Automatic post archiving
- Clean separation of concerns via **Hexagonal Architecture**
- Session handling with cookies

This is a learning-oriented project covering:
- REST APIs
- Authentication & cookies
- S3 integration
- PostgreSQL with SQL
- Logging with `log/slog`
- Unit testing
- Basic frontend integration
- Concurrency fundamentals

---

## 🚀 Features

- **Anonymous posting** — no account required.
- **Images** — uploaded to S3-compatible storage.
- **Rick & Morty avatars** — assigned per user session.
- **Session tracking** — via secure cookies (1 week lifetime).
- **Post lifecycle**:
  - Without comments → auto-delete from catalog after 10 min.
  - With comments → auto-delete 15 min after last comment.
- **Archiving** — deleted posts accessible in archive (read-only).
- **Replies** — to posts and specific comments.
- **Hexagonal Architecture** — clean separation of domain, infrastructure, and UI.

---

## 🏗 Architecture

### Layers
1. **Domain Layer (Core Logic)**  
   - Business rules for posts, comments, and sessions.
   - Defines interfaces for storage and external services.

2. **Infrastructure Layer (Adapters)**  
   - PostgreSQL adapter (data persistence)
   - S3 adapter (image storage)
   - Rick & Morty API adapter (avatars)

3. **Interface Layer**  
   - HTTP handlers for REST API
   - Middleware for authentication/session management
   - HTML templates for UI

---

## 📂 Templates Provided
- `catalog.html` — list of active posts
- `archive.html` — archived posts
- `post.html` — single post + comments
- `archive-post.html` — archived post view
- `create-post.html` — new post form
- `error.html` — error page

---

## 💾 Data Storage
- **PostgreSQL**: posts, comments, sessions, metadata.
- **S3-Compatible Storage**: images (MinIO recommended for local dev).
- **Avatars**: retrieved dynamically (not stored locally).

---

## 🛠 Installation & Setup

### 1. Clone the repo
```sh
git clone https://github.com/yourusername/1337b04rd.git
cd 1337b04rd
