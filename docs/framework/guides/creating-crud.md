---
title: "Creating a CRUD Application"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Creating a CRUD Application

## Abstract

A practical guide to building a complete CRUD application — from
model to routes to views — using the MVC + Services pattern.

## Table of Contents

1. [Overview](#1-overview)
2. [Create Model](#2-create-model)
3. [Run Migration](#3-run-migration)
4. [Seed Sample Data](#4-seed-sample-data)
5. [Create Service](#5-create-service)
6. [Create Controller](#6-create-controller)
7. [Register Routes](#7-register-routes)
8. [Create Views](#8-create-views)
9. [Add Validation](#9-add-validation)
10. [Add Flash Messages](#10-add-flash-messages)
11. [References](#11-references)

## 1. Overview

We will build a `Post` resource with index, create, show, edit, and
delete functionality.

## 2. Create Model

```bash
framework make:model Post
```

Edit `database/models/post.go`:

```go
package models

type Post struct {
    BaseModel
    Title   string `gorm:"size:200;not null" json:"title"`
    Body    string `gorm:"type:text;not null" json:"body"`
    UserID  uint   `gorm:"not null" json:"user_id"`
    User    User   `gorm:"foreignKey:UserID" json:"user"`
}
```

## 3. Run Migration

```bash
framework migrate
```

## 4. Seed Sample Data

In `database/seeders/posts.go`:

```go
package seeders

import (
    "yourframework/database/models"
    "gorm.io/gorm"
)

func SeedPosts(db *gorm.DB) {
    posts := []models.Post{
        {Title: "First Post", Body: "Hello, world!", UserID: 1},
        {Title: "Second Post", Body: "Another post.", UserID: 1},
    }
    for _, p := range posts {
        db.FirstOrCreate(&p, models.Post{Title: p.Title})
    }
}
```

```bash
framework db:seed
```

## 5. Create Service

```bash
framework make:service PostService
```

Edit `app/services/postservice.go`:

```go
package services

import (
    "yourframework/database/models"
    "gorm.io/gorm"
)

type PostService struct {
    DB *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
    return &PostService{DB: db}
}

func (s *PostService) GetAll(page, perPage int) ([]models.Post, int64, error) {
    var posts []models.Post
    var total int64
    s.DB.Model(&models.Post{}).Count(&total)
    err := s.DB.Preload("User").
        Offset((page - 1) * perPage).
        Limit(perPage).
        Order("created_at DESC").
        Find(&posts).Error
    return posts, total, err
}

func (s *PostService) GetByID(id uint) (*models.Post, error) {
    var post models.Post
    err := s.DB.Preload("User").First(&post, id).Error
    return &post, err
}

func (s *PostService) Create(post *models.Post) error {
    return s.DB.Create(post).Error
}

func (s *PostService) Update(post *models.Post) error {
    return s.DB.Save(post).Error
}

func (s *PostService) Delete(id uint) error {
    return s.DB.Delete(&models.Post{}, id).Error
}
```

## 6. Create Controller

```bash
framework make:controller PostController
```

Edit to use `PostService`:

```go
package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "yourframework/app/services"
    "yourframework/database/models"
)

type PostController struct {
    Service *services.PostService
}

func (ctrl *PostController) Index(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    posts, total, _ := ctrl.Service.GetAll(page, 10)
    c.HTML(http.StatusOK, "posts/index.html", gin.H{
        "posts": posts,
        "total": total,
        "page":  page,
    })
}

func (ctrl *PostController) Show(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    post, err := ctrl.Service.GetByID(uint(id))
    if err != nil {
        c.HTML(http.StatusNotFound, "errors/404.html", nil)
        return
    }
    c.HTML(http.StatusOK, "posts/show.html", gin.H{"post": post})
}

func (ctrl *PostController) Store(c *gin.Context) {
    var post models.Post
    post.Title = c.PostForm("title")
    post.Body = c.PostForm("body")
    post.UserID = 1 // from auth in production

    if err := ctrl.Service.Create(&post); err != nil {
        c.Redirect(http.StatusFound, "/posts/create")
        return
    }
    c.Redirect(http.StatusFound, "/posts")
}

func (ctrl *PostController) Update(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    post, _ := ctrl.Service.GetByID(uint(id))
    post.Title = c.PostForm("title")
    post.Body = c.PostForm("body")
    ctrl.Service.Update(post)
    c.Redirect(http.StatusFound, "/posts")
}

func (ctrl *PostController) Destroy(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    ctrl.Service.Delete(uint(id))
    c.Redirect(http.StatusFound, "/posts")
}
```

## 7. Register Routes

In `routes/web.go`:

```go
postCtrl := &controllers.PostController{
    Service: services.NewPostService(db),
}
router.Resource("/posts", postCtrl)
```

This creates:

| Method | Path | Action |
|--------|------|--------|
| GET | `/posts` | `Index` |
| GET | `/posts/:id` | `Show` |
| POST | `/posts` | `Store` |
| PUT | `/posts/:id` | `Update` |
| DELETE | `/posts/:id` | `Destroy` |

## 8. Create Views

Create templates in `resources/views/posts/`:

### `index.html`

```html
{{template "layouts/app" .}}
{{define "content"}}
<h1>Posts</h1>
<a href="/posts/create">New Post</a>
{{range .posts}}
<div>
    <h2><a href="/posts/{{.ID}}">{{.Title}}</a></h2>
    <p>by {{.User.Name}} — {{.CreatedAt.Format "Jan 2, 2006"}}</p>
</div>
{{end}}
{{end}}
```

### `show.html`

```html
{{template "layouts/app" .}}
{{define "content"}}
<h1>{{.post.Title}}</h1>
<p>{{.post.Body}}</p>
<a href="/posts">Back</a>
{{end}}
```

## 9. Add Validation

Use the built-in validator before creating:

```go
func (ctrl *PostController) Store(c *gin.Context) {
    v := validation.New()
    v.Required("title", c.PostForm("title"))
    v.MinLength("title", c.PostForm("title"), 3)
    v.Required("body", c.PostForm("body"))

    if v.HasErrors() {
        session.FlashErrors(c, v.Errors)
        session.FlashOldInput(c, map[string]string{
            "title": c.PostForm("title"),
            "body":  c.PostForm("body"),
        })
        c.Redirect(http.StatusFound, "/posts/create")
        return
    }
    // ... create post
}
```

## 10. Add Flash Messages

Display success/error messages in templates:

```go
// In controller after successful create:
session.Flash(c, "success", "Post created successfully!")
c.Redirect(http.StatusFound, "/posts")
```

```html
<!-- In template -->
{{if .success}}
<div class="alert alert-success">{{.success}}</div>
{{end}}
{{if .errors}}
<ul class="alert alert-danger">
    {{range $field, $msg := .errors}}
    <li>{{$field}}: {{$msg}}</li>
    {{end}}
</ul>
{{end}}
```

## 11. References

- [Models](../data/models.md)
- [Services Layer](../infrastructure/services-layer.md)
- [Controllers](../http/controllers.md)
- [Routing](../http/routing.md)
- [Views](../http/views.md)
- [Requests & Validation](../http/requests-validation.md)
- [Sessions — Flash Messages](../security/sessions.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
