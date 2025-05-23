Here is the content for a file named `technical_design.md` containing the technical design choices for the `squirrel` project:

---

```markdown
# 🛠️ Squirrel - Technical Design Specification (v0.1)

**Project Name:** `squirrel`  
**Owner:** Karim Fanous  
**Stage:** Technical Design Spec  
**Date:** 2025-05-16  
**Version:** v0.1  
**Stack:** Go (backend) + React (frontend) + PostgreSQL (DB) + GraphQL API

---

## 🧱 1. System Architecture Overview

```

+---------------------------+       +---------------------+
\|     React Web Frontend    | <---> |    GraphQL Gateway   |
+---------------------------+       +---------------------+
|
|
+--------------------+
\|     Go Backend      |
\| - Core logic        |
\| - Parser (HTML, Email) |
\| - Ingestion          |
+--------------------+
|
|
+---------------------+
\|     PostgreSQL DB    |
+---------------------+

````

---

## 🧰 2. Technology Stack

| Layer            | Stack / Tools                                          |
|------------------|--------------------------------------------------------|
| Frontend         | React + TypeScript + Vite + Tailwind CSS               |
| Routing          | React Router                                           |
| State Management | React Query (for GraphQL hooks)                        |
| GraphQL API      | 99designs/gqlgen (Go)                                  |
| Backend Language | Go (Golang)                                            |
| DB Layer         | sqlc or gorm                                           |
| DB Engine        | PostgreSQL                                             |
| Email Ingest     | IMAP/SMTP via `emersion/go-imap`, `mhale/smtpd`       |
| Parsing HTML     | `mauidb/go-readability` or `PuerkitoBio/goquery`       |
| Full-text Search | PostgreSQL `tsvector` or `pgroonga` (optional)         |
| Deployment       | Docker, Make, fly.io/localhost                         |

---

## 📦 3. Backend Components

### Entry Service
- Handles create/read/tag/search logic
- GraphQL resolvers

### Parser Service
- Article content extraction
- Title/metadata enhancement

### Email Ingestion Service
- IMAP polling or SMTP receiver
- Converts email to `task` entry

### Task Scheduler (Optional)
- Periodic email polling
- Background parsing jobs

---

## 🧾 4. Data Model (PostgreSQL)

### `entries` Table

```sql
CREATE TABLE entries (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  url TEXT,
  content TEXT,
  tags TEXT[],
  entry_type TEXT NOT NULL CHECK (entry_type IN ('note', 'article', 'task')),
  read BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP
);
````

### Indexes

```sql
CREATE INDEX idx_entry_type ON entries(entry_type);
CREATE INDEX idx_tags ON entries USING GIN(tags);
CREATE INDEX idx_search ON entries USING GIN(to_tsvector('english', content));
```

---

## 🧩 5. GraphQL Schema (Sample)

```graphql
type Entry {
  id: ID!
  title: String!
  url: String
  content: String!
  tags: [String!]!
  entryType: String!
  read: Boolean!
  createdAt: String!
  updatedAt: String
}

type Query {
  getEntry(id: ID!): Entry
  listEntries(entryType: String, tag: String): [Entry!]!
  searchEntries(query: String!): [Entry!]!
}

type Mutation {
  addEntry(input: AddEntryInput!): Entry!
  markRead(id: ID!): Entry!
  tagEntry(id: ID!, tags: [String!]!): Entry!
}
```

---

## 💻 6. Frontend Design (React)

### Key Views

* Entry List View
* Entry Detail View
* Add Entry Modal

### UI Stack

* Tailwind CSS
* React Router
* GraphQL Codegen
* React Query

---

## 🔐 7. Auth & Security (Future-ready)

* MVP: local-only, no auth
* Future: JWT or basic auth, CLI token access

---

## 🧪 8. Testing Plan

* Unit tests: parser, entry service
* Resolver tests
* CLI integration tests
* IMAP inbox simulation
* E2E tests with Playwright/Cypress

---

## 🚀 9. Dev Environment

* Docker Compose for Postgres
* `.env` config file
* `make dev` script
* Hot reload via `air` (Go)

---

## 📈 10. Stretch Goals

* Markdown notes
* Task reminders and due dates
* Email threads as comments
* Obsidian export
* Mobile PWA

```

---

Would you like me to now scaffold the actual directory structure (`/cmd`, `/internal`, etc.) or generate starter code for one of the components (e.g., GraphQL schema + resolver)?
```
