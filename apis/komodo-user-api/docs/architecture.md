# Architecture: User Management API

## 1. Overview
The Komodo User API is a RESTful service responsible for handling user identities, authentication, profile management, and preferences. It serves as the primary authority for user data across the platform.

### Core Responsibilities
* **Identity Management:** Registration, updates, and account deletion.
* **Authentication:** Login, logout, and password resets.
* **Authorization:** Role-based access control (RBAC).

---

## 2. System Design
The service follows a **Layered Architecture** to ensure separation of concerns and testability.

| Layer | Responsibility |
| :--- | :--- |
| **Controller** | Handles HTTP requests, input validation, and status codes. |
| **Service** | Contains business logic (e.g., "Can this user upgrade their plan?"). |
| **Repository** | Manages data persistence and database queries. |
| **Model** | Defines the data structures and database schema. |

---

## 3. Technology Stack
* **Language:** Go 1.26.0
* **Framework:** Go HTTP (Mux)
* **Database:** DynamoDB (AWS)
* **Auth:** OAuth 2.0 (as JWT)
* **Documentation:** Swagger/OpenAPI 3.0

---

## 4. Data Flow
### Authentication Flow
1.  **Client** sends credentials to `/login`.
2.  **API** validates credentials against the **Database**.
3.  On success, the API generates a **JWT** and stores a refresh token in **Redis**.
4.  **Client** includes the JWT in the `Authorization: Bearer <token>` header for subsequent requests.

---

## 5. Database Schema
We utilize a NoSQL schema to maintain data integrity.

> **Note:** Passwords must never be stored in plain text. We use `bcrypt` with a salt factor of 12.

### Key Tables
* **`users`**: `id`, `email`, `password_hash`, `created_at`
* **`profiles`**: `user_id`, `first_name`, `last_name`, `avatar_url`
* **`roles`**: `id`, `role_name` (Admin, Editor, Viewer)

---

## 6. Security Considerations
* **Rate Limiting:** Maximum 100 requests per minute per IP.
* **CORS:** Restricted to trusted domains only.
* **Encryption:** All data in transit is encrypted via TLS 1.3.
* **PII:** Personally Identifiable Information is encrypted at rest using AES-256.

---

## 7. Scalability
* **Horizontal Scaling:** The service is stateless and can be deployed across multiple containers using Kubernetes.
* **Read Replicas:** Database read operations are distributed across replicas to reduce load on the primary instance.