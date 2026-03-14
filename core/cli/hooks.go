package cli

import (
	"github.com/raiworks/rapidgo/v2/core/app"
	"github.com/raiworks/rapidgo/v2/core/container"
	"github.com/raiworks/rapidgo/v2/core/router"
	"github.com/raiworks/rapidgo/v2/core/scheduler"
	"github.com/raiworks/rapidgo/v2/core/service"
	"gorm.io/gorm"
)

// BootstrapFunc registers service providers on the application for the given mode.
type BootstrapFunc func(a *app.App, mode service.Mode)

// RouteRegistrar registers routes on the router for a given mode.
type RouteRegistrar func(r *router.Router, c *container.Container, mode service.Mode)

// JobRegistrar registers application job handlers with the queue dispatcher.
type JobRegistrar func()

// ScheduleRegistrar registers scheduled tasks on the scheduler.
type ScheduleRegistrar func(s *scheduler.Scheduler, a *app.App)

// ModelRegistryFunc returns all model structs for AutoMigrate.
type ModelRegistryFunc func() []interface{}

// SeederFunc runs database seeders. If name is empty, runs all seeders.
type SeederFunc func(db *gorm.DB, name string) error

// SeederListFunc returns a list of registered seeder names.
type SeederListFunc func() []string

var (
	bootstrapFn       BootstrapFunc
	routeRegistrar    RouteRegistrar
	jobRegistrar      JobRegistrar
	scheduleRegistrar ScheduleRegistrar
	modelRegistryFn   ModelRegistryFunc
	seederFn          SeederFunc
	seederListFn      SeederListFunc
)

// SetBootstrap sets the function that registers service providers during app initialization.
func SetBootstrap(fn BootstrapFunc) { bootstrapFn = fn }

// SetRoutes sets the function that registers routes on the router.
func SetRoutes(fn RouteRegistrar) { routeRegistrar = fn }

// SetJobRegistrar sets the function that registers background job handlers.
func SetJobRegistrar(fn JobRegistrar) { jobRegistrar = fn }

// SetScheduleRegistrar sets the function that registers scheduled tasks.
func SetScheduleRegistrar(fn ScheduleRegistrar) { scheduleRegistrar = fn }

// SetModelRegistry sets the function that returns all model structs for AutoMigrate.
func SetModelRegistry(fn ModelRegistryFunc) { modelRegistryFn = fn }

// SetSeeder sets the function that runs database seeders.
func SetSeeder(fn SeederFunc) { seederFn = fn }

// SetSeederList sets the function that returns available seeder names for db:seed --list.
func SetSeederList(fn SeederListFunc) { seederListFn = fn }
