# 📦 go-db-tx

**go-db-tx** is a lightweight, reusable transaction management library for Go that provides clean, context-based transaction control for **PostgreSQL (`database/sql`)** and **TimescaleDB (`pgx/pgxpool`)**.

It is built with **Clean Architecture** and **usecase-driven design** in mind, enabling seamless transaction propagation across repository layers without tightly coupling business logic to database implementations.

---

## ✨ Key Features

* Context-based transaction handling
* Supports PostgreSQL and TimescaleDB transactions
* Automatic commit and rollback management
* Panic-safe transaction recovery
* Clean Architecture & DDD friendly
* Eliminates repetitive transaction boilerplate across projects

---

## 🎯 Use Cases

Ideal for Go backend services where:

* Multiple repositories need to share the same transaction
* Business logic resides in the usecase layer
* You want consistent and centralized transaction handling
* You want to avoid rewriting transaction control code in every project

---

## 🚀 Why go-db-tx?

Managing database transactions across complex service layers often leads to duplicated and error-prone code. **go-db-tx** abstracts this concern into a reusable, composable component, allowing developers to focus on business logic while maintaining transactional consistency and safety.

---

## 🧱 Architecture Fit

This library fits naturally into projects following:

* Clean Architecture
* Domain-Driven Design (DDD)
* Hexagonal / Onion Architecture

It works by propagating transactions through `context.Context`, making it easy to integrate with existing repository and usecase layers.

---

## 📄 License

MIT License
