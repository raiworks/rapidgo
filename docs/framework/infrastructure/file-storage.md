---
title: "File Storage"
version: "0.1.0"
status: "Final"
date: "2026-03-05"
last_updated: "2026-03-10"
authors:
  - "RAiWorks"
supersedes: ""
---

# File Storage

## Abstract

This document covers the driver-based file storage system — the
storage interface, local disk and S3 drivers, configuration, and
the upload controller pattern.

## Table of Contents

1. [Terminology](#1-terminology)
2. [Configuration](#2-configuration)
3. [Storage Interface](#3-storage-interface)
4. [Local Driver](#4-local-driver)
5. [S3 Driver](#5-s3-driver)
6. [Upload Controller](#6-upload-controller)
7. [Security Considerations](#7-security-considerations)
8. [References](#8-references)

## 1. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
document are to be interpreted as described in [RFC 2119].

- **Storage driver** — A backend implementing the `Driver` interface
  for file operations.

## 2. Configuration

`.env`:

```env
STORAGE_DRIVER=local
STORAGE_LOCAL_PATH=storage/uploads

# S3 driver only
AWS_REGION=us-east-1
AWS_BUCKET=my-bucket
AWS_ACCESS_KEY_ID=xxx
AWS_SECRET_ACCESS_KEY=xxx
```

## 3. Storage Interface

```go
package storage

import "io"

type Driver interface {
    Put(path string, content io.Reader) (string, error)
    Get(path string) (io.ReadCloser, error)
    Delete(path string) error
    URL(path string) string
}
```

## 4. Local Driver

Stores files on the local filesystem:

```go
type LocalDriver struct {
    BasePath string // e.g. "storage/uploads"
    BaseURL  string // e.g. "/static/uploads"
}

func (d *LocalDriver) Put(path string, content io.Reader) (string, error) {
    fullPath := filepath.Join(d.BasePath, path)
    if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
        return "", err
    }
    f, err := os.Create(fullPath)
    if err != nil {
        return "", err
    }
    defer f.Close()
    if _, err := io.Copy(f, content); err != nil {
        return "", err
    }
    return path, nil
}

func (d *LocalDriver) Get(path string) (io.ReadCloser, error) {
    return os.Open(filepath.Join(d.BasePath, path))
}

func (d *LocalDriver) Delete(path string) error {
    return os.Remove(filepath.Join(d.BasePath, path))
}

func (d *LocalDriver) URL(path string) string {
    return d.BaseURL + "/" + path
}
```

## 5. S3 Driver

Library: `github.com/aws/aws-sdk-go-v2`

```go
type S3Driver struct {
    Client *s3.Client
    Bucket string
    Region string
}

func (d *S3Driver) Put(path string, content io.Reader) (string, error) {
    _, err := d.Client.PutObject(context.Background(), &s3.PutObjectInput{
        Bucket: aws.String(d.Bucket),
        Key:    aws.String(path),
        Body:   content,
    })
    return path, err
}

func (d *S3Driver) Get(path string) (io.ReadCloser, error) {
    out, err := d.Client.GetObject(context.Background(), &s3.GetObjectInput{
        Bucket: aws.String(d.Bucket),
        Key:    aws.String(path),
    })
    if err != nil {
        return nil, err
    }
    return out.Body, nil
}

func (d *S3Driver) Delete(path string) error {
    _, err := d.Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
        Bucket: aws.String(d.Bucket),
        Key:    aws.String(path),
    })
    return err
}

func (d *S3Driver) URL(path string) string {
    return "https://" + d.Bucket + ".s3." + d.Region + ".amazonaws.com/" + path
}
```

## 6. Upload Controller

```go
func Upload(c *gin.Context) {
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        responses.Error(c, 400, "no file uploaded")
        return
    }
    defer file.Close()

    path := "uploads/" + helpers.RandomString(16) + filepath.Ext(header.Filename)
    savedPath, err := storageDriver.Put(path, file)
    if err != nil {
        responses.Error(c, 500, "failed to save file")
        return
    }
    responses.Success(c, gin.H{
        "path": savedPath,
        "url":  storageDriver.URL(savedPath),
    })
}
```

## 7. Security Considerations

- AWS credentials **MUST** be stored in `.env` and **MUST NOT** be
  committed to version control.
- File paths **MUST** be sanitized — never use the original filename
  directly. Generate random filenames.
- Validate file types and sizes before storage.
- The upload directory **MUST NOT** allow execution of uploaded files
  (e.g., no `.go`, `.sh` execution).
- S3 buckets **SHOULD** use private ACLs with pre-signed URLs for
  access control.
- Local uploads served via `r.Static()` are publicly accessible —
  use access control middleware for private files.

## 8. References

- [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)
- [Helpers — RandomString](../reference/helpers-reference.md)
- [Configuration](../core/configuration.md)
- [Views — Static File Serving](../http/views.md#9-static-file-serving)

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1.0 | 2026-03-05 | RAiWorks | Initial draft |
