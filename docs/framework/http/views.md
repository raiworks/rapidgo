---
title: "Views & Templates"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# Views & Templates

## Abstract

This document covers the server-side rendering (SSR) system using
Go's `html/template` package, template directory conventions, data
passing, template functions (including named route helpers), the Templ
alternative, and static file serving.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Template Engine Setup](#2-template-engine-setup)
3. [Template Directory Structure](#3-template-directory-structure)
4. [Passing Data to Templates](#4-passing-data-to-templates)
5. [Template Functions](#5-template-functions)
6. [Form Templates with Validation](#6-form-templates-with-validation)
7. [Flash Messages in Templates](#7-flash-messages-in-templates)
8. [Templ (Alternative)](#8-templ-alternative)
9. [Static File Serving](#9-static-file-serving)
10. [Security Considerations](#10-security-considerations)
11. [References](#11-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **SSR** — Server-Side Rendering; HTML is generated on the server
  and sent to the client.
- **Template function** — A Go function exposed to templates via
  `FuncMap`.

## 2. Template Engine Setup

The framework uses Go's standard `html/template` package, loaded via
Gin's `LoadHTMLGlob`:

```go
func SetupRouter() *gin.Engine {
    r := gin.Default()
    r.LoadHTMLGlob("resources/views/**/*")
    // ...
    return r
}
```

Templates are stored in `resources/views/`.

## 3. Template Directory Structure

```text
resources/
└── views/
    ├── layouts/
    │   └── base.html
    ├── home.html
    ├── users/
    │   ├── index.html
    │   ├── show.html
    │   ├── create.html
    │   └── edit.html
    └── posts/
        ├── index.html
        ├── show.html
        ├── create.html
        └── edit.html
```

## 4. Passing Data to Templates

Use `gin.H` (alias for `map[string]interface{}`) to pass data:

```go
func Home(c *gin.Context) {
    c.HTML(http.StatusOK, "home.html", gin.H{
        "title": "Welcome",
    })
}

func ShowUser(c *gin.Context) {
    c.HTML(http.StatusOK, "users/show.html", gin.H{
        "title": "User Profile",
        "user":  user,
    })
}
```

Access in templates:

```html
<h1>{{.title}}</h1>
<p>Name: {{.user.Name}}</p>
<p>Email: {{.user.Email}}</p>
```

## 5. Template Functions

Register `Route` as a template function for URL generation:

```go
r.SetFuncMap(template.FuncMap{
    "route": router.Route,
})
r.LoadHTMLGlob("resources/views/**/*")
```

Use in templates:

```html
<a href="{{route `users.show` .User.ID}}">View Profile</a>
<a href="{{route `posts.edit` .Post.ID}}">Edit Post</a>
<a href="{{route `home`}}">Home</a>
```

## 6. Form Templates with Validation

Forms **MUST** include a CSRF token. Use old input and error display
for SSR validation flows:

```html
<form method="POST" action="/users">
    <input type="hidden" name="_csrf_token" value="{{.csrf_token}}">

    <label>Name</label>
    <input type="text" name="name" value="{{.old.name}}">
    {{if .errors.name}}<span class="error">{{index .errors.name 0}}</span>{{end}}

    <label>Email</label>
    <input type="email" name="email" value="{{.old.email}}">
    {{if .errors.email}}<span class="error">{{index .errors.email 0}}</span>{{end}}

    <label>Password</label>
    <input type="password" name="password">
    {{if .errors.password}}<span class="error">{{index .errors.password 0}}</span>{{end}}

    <label>Confirm Password</label>
    <input type="password" name="password_confirmation">

    <button type="submit">Create</button>
</form>
```

Controller passes old input and errors on validation failure:

```go
if !v.Valid() {
    c.HTML(http.StatusUnprocessableEntity, "users/create.html", gin.H{
        "errors": v.Errors(),
        "old":    gin.H{"name": name, "email": email},
    })
    return
}
```

## 7. Flash Messages in Templates

Display one-time flash messages passed via session:

```html
{{if .success}}
<div class="alert alert-success">{{.success}}</div>
{{end}}

{{if .errors}}
<div class="alert alert-danger">
    <ul>
    {{range $field, $msgs := .errors}}
        {{range $msgs}}<li>{{.}}</li>{{end}}
    {{end}}
    </ul>
</div>
{{end}}
```

See [Sessions](../security/sessions.md) for flash message API details.

## 8. Templ (Alternative)

For a more modern, type-safe approach, consider **Templ**
(`github.com/a-h/templ`). Templ provides compiled, type-safe templates
that catch errors at build time rather than runtime.

Key benefits:
- **Type safety** — templates are compiled Go code
- **IDE support** — autocompletion and error checking
- **Performance** — compiled templates, no parsing overhead

Templ is **OPTIONAL**; `html/template` remains the default.

## 9. Static File Serving

### Static Assets

Serve CSS, JS, images from `resources/static/`:

```go
r.Static("/static", "./resources/static")
```

### User Uploads

Serve uploaded files from `storage/uploads/`:

```go
r.Static("/uploads", "./storage/uploads")
```

### Favicon

```go
r.StaticFile("/favicon.ico", "./resources/static/favicon.ico")
```

### Complete Setup

```go
func SetupRouter() *gin.Engine {
    r := gin.Default()

    // Static assets (CSS, JS, images)
    r.Static("/static", "./resources/static")

    // User uploads
    r.Static("/uploads", "./storage/uploads")

    // Favicon
    r.StaticFile("/favicon.ico", "./resources/static/favicon.ico")

    r.LoadHTMLGlob("resources/views/**/*")

    // ... routes
    return r
}
```

Reference in templates:

```html
<link rel="stylesheet" href="/static/css/app.css">
<script src="/static/js/app.js"></script>
<img src="/uploads/{{.user.Avatar}}" alt="Avatar">
```

## 10. Security Considerations

- Go's `html/template` **auto-escapes** HTML output by default,
  protecting against XSS.
- Forms **MUST** include CSRF tokens via hidden fields.
- Uploaded file paths **MUST** be sanitized — never allow user input
  to dictate filesystem paths directly.
- Static file directories **SHOULD** be limited to specific paths;
  avoid serving the entire project root.

## 11. References

- [Controllers](controllers.md)
- [Routing — Named Routes](routing.md#6-named-routes--url-generation)
- [Requests & Validation](requests-validation.md)
- [Sessions — Flash Messages](../security/sessions.md)
- [CSRF Protection](../security/csrf.md)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
