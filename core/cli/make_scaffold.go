package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

// scaffold generates a Go source file from a template.
// It prevents overwriting existing files and creates directories as needed.
func scaffold(kind, name, tpl, dir string, out io.Writer) error {
	filename := toSnakeCase(name) + ".go"
	path := filepath.Join(dir, filename)

	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("file already exists or cannot be created: %w", err)
	}
	defer f.Close()

	t := template.Must(template.New(kind).Parse(tpl))
	if err := t.Execute(f, map[string]string{"Name": name}); err != nil {
		return fmt.Errorf("failed to write template: %w", err)
	}

	fmt.Fprintf(out, "%s created: %s\n", kind, path)
	return nil
}

// --- Commands ---

var makeControllerCmd = &cobra.Command{
	Use:   "make:controller [name]",
	Short: "Create a new controller",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return scaffold("Controller", args[0], controllerTpl, "http/controllers", cmd.OutOrStdout())
	},
}

var makeModelCmd = &cobra.Command{
	Use:   "make:model [name]",
	Short: "Create a new model",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return scaffold("Model", args[0], modelTpl, "database/models", cmd.OutOrStdout())
	},
}

var makeServiceCmd = &cobra.Command{
	Use:   "make:service [name]",
	Short: "Create a new service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return scaffold("Service", args[0], serviceTpl, "app/services", cmd.OutOrStdout())
	},
}

var makeProviderCmd = &cobra.Command{
	Use:   "make:provider [name]",
	Short: "Create a new service provider",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return scaffold("Provider", args[0], providerTpl, "app/providers", cmd.OutOrStdout())
	},
}

// --- Templates ---

var controllerTpl = `package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type {{.Name}} struct{}

func (ctrl *{{.Name}}) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "{{.Name}} index"})
}

func (ctrl *{{.Name}}) Show(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (ctrl *{{.Name}}) Store(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (ctrl *{{.Name}}) Update(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "updated"})
}

func (ctrl *{{.Name}}) Destroy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
`

var modelTpl = `package models

type {{.Name}} struct {
	BaseModel
	// Add fields here
}
`

var serviceTpl = `package services

import "gorm.io/gorm"

type {{.Name}} struct {
	DB *gorm.DB
}

func New{{.Name}}(db *gorm.DB) *{{.Name}} {
	return &{{.Name}}{DB: db}
}

// Add service methods here
`

var providerTpl = `package providers

import "github.com/raiworks/rapidgo/v2/core/container"

type {{.Name}} struct{}

func (p *{{.Name}}) Register(c *container.Container) {
	// Bind services into the container
}

func (p *{{.Name}}) Boot(c *container.Container) {
	// Run after all providers are registered
}
`

var makeSeederCmd = &cobra.Command{
	Use:   "make:seeder [name]",
	Short: "Create a new database seeder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return scaffold("Seeder", args[0], seederTpl, "database/seeders", cmd.OutOrStdout())
	},
}

var seederTpl = `package seeders

import "gorm.io/gorm"

func {{.Name}}Seeder(db *gorm.DB) error {
	// Seed data here
	return nil
}
`

var makeModuleCmd = &cobra.Command{
	Use:   "make:module [name]",
	Short: "Create a new domain module with models, service, controller, and routes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		dir := filepath.Join("modules", toSnakeCase(name))
		out := cmd.OutOrStdout()

		files := []struct {
			kind string
			tpl  string
			file string
		}{
			{"Models", moduleModelsTpl, "models.go"},
			{"Service", moduleServiceTpl, "service.go"},
			{"Controller", moduleControllerTpl, "controller.go"},
			{"Routes", moduleRoutesTpl, "routes.go"},
		}

		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("failed to create module directory: %w", err)
		}

		for _, f := range files {
			path := filepath.Join(dir, f.file)
			file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("file already exists or cannot be created: %w", err)
			}
			t := template.Must(template.New(f.kind).Parse(f.tpl))
			if err := t.Execute(file, map[string]string{
				"Name":    name,
				"Package": toSnakeCase(name),
			}); err != nil {
				file.Close()
				return fmt.Errorf("failed to write template: %w", err)
			}
			file.Close()
			fmt.Fprintf(out, "%s created: %s\n", f.kind, path)
		}
		return nil
	},
}

var moduleModelsTpl = `package {{.Package}}

import "github.com/raiworks/rapidgo/v2/database/models"

type {{.Name}} struct {
	models.BaseModel
	// Add fields here
}
`

var moduleServiceTpl = `package {{.Package}}

import "gorm.io/gorm"

type {{.Name}}Service struct {
	DB *gorm.DB
}

func New{{.Name}}Service(db *gorm.DB) *{{.Name}}Service {
	return &{{.Name}}Service{DB: db}
}

// Add service methods here
`

var moduleControllerTpl = `package {{.Package}}

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type {{.Name}}Controller struct {
	Service *{{.Name}}Service
}

func (ctrl *{{.Name}}Controller) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "{{.Name}} index"})
}
`

var moduleRoutesTpl = `package {{.Package}}

import "github.com/gin-gonic/gin"

// RegisterRoutes adds this module's routes to the given router group.
func RegisterRoutes(group *gin.RouterGroup, ctrl *{{.Name}}Controller) {
	group.GET("/", ctrl.Index)
}
`
