# 🐿️ Squirrel - Functional Specification (v0.1)

**Project Name:** `squirrel`  
**Owner:** Karim Fanous  
**Stage:** Functional Spec  
**Date:** 2025-05-16  
**Version:** v0.1

---

## 🧭 Purpose

`squirrel` is a privacy-first, self-hosted tool that allows users to capture **articles**, **notes**, and **tasks** via CLI, browser, and email. It supports local-first usage, tagging, and search — all served from a minimal Go binary.

---

## 🎯 Goals & Objectives

- ✅ Save and parse articles by URL
- ✅ Capture freeform notes (todos, ideas, reminders)
- ✅ Create tasks by sending an email
- ✅ Tag, read, and search entries
- ✅ Offer CLI and web interfaces

---

## 👥 Target User

- Developers, knowledge workers, and privacy-minded users
- People who email themselves notes or tasks
- Users who want a unified inbox for articles, notes, and todos

---

## 📦 Features & Requirements

### 1. Entry Types

| Type    | Description                            |
|---------|----------------------------------------|
| Article | Saved via URL; parsed content fetched  |
| Note    | Text saved via CLI or UI               |
| Task    | Notes submitted via **email**, parsed into entries |

---

### 2. CLI Interface

```bash
squirrel add <url or note> --tag <tag1,tag2>
squirrel list [--tag] [--type article|note|task]
squirrel read <id>
squirrel serve
