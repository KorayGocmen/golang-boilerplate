package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/aws"
	"github.com/koraygocmen/golang-boilerplate/internal/config"
	"github.com/koraygocmen/golang-boilerplate/internal/database"
	"github.com/koraygocmen/golang-boilerplate/internal/env"
	"github.com/koraygocmen/golang-boilerplate/internal/errhandle"
	"github.com/koraygocmen/golang-boilerplate/internal/logger"
	"github.com/koraygocmen/golang-boilerplate/internal/transport/handler"
	"github.com/koraygocmen/golang-boilerplate/internal/transport/middleware"
	"github.com/koraygocmen/golang-boilerplate/internal/transport/router"
	_ "github.com/koraygocmen/golang-boilerplate/migrations"
	"github.com/koraygocmen/golang-boilerplate/pkg/duration"
	_ "github.com/koraygocmen/golang-boilerplate/seeds"
	"github.com/spf13/cobra"
)

var (
	SHASUM string

	cmdRoot = &cobra.Command{}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Ping command.
	{
		cmdRoot.AddCommand(&cobra.Command{
			Use:   "ping",
			Short: "Ping application bin",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("pong")
				os.Exit(0)
			},
		})
	}

	// Version command.
	{
		cmdRoot.AddCommand(&cobra.Command{
			Use:   "version",
			Short: "Print the version number of application bin",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(SHASUM)
				os.Exit(0)
			},
		})
	}

	// Setting up the database commands.
	// The root command is the db command.
	{
		cmdDB := &cobra.Command{
			Use:   "db",
			Short: "Database operations",
			Args:  cobra.MinimumNArgs(1),
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				// Load env variables.
				env.Init()

				// Initialize aws client.
				if err := aws.Init(); err != nil {
					err = fmt.Errorf("aws init error: %w", err)
					log.Fatal(err)
				}

				// Load config.
				config.Load()

				// Set up logger.
				logWriter, err := logger.New(logger.Config{
					Level:  config.Log.Level,
					Mode:   config.Log.Mode,
					SHASUM: SHASUM,

					Syslog: config.Log.Syslog,
					File:   config.Log.File,
					AWS:    config.Log.AWS,
				})
				if err != nil {
					err = fmt.Errorf("logger new error: %w", err)
					log.Fatal(err)
				}

				// Set the logger.
				logger.Logger = logWriter

				// Initialize database.
				if database.DB, err = database.Connect(logger.Logger, config.Database); err != nil {
					err = fmt.Errorf("database new error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			},
		}
		cmdRoot.AddCommand(cmdDB)

		// db list command.
		cmdDB.AddCommand(&cobra.Command{
			Use:   "list",
			Short: "db list",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				if err := database.DB.List(); err != nil {
					err = fmt.Errorf("db list error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			},
		})

		// db up command.
		cmdDB.AddCommand(&cobra.Command{
			Use:   "up",
			Short: "db up",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				if err := database.DB.Up(ctx); err != nil {
					err = fmt.Errorf("db up error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			},
		})

		// db up-by-one command.
		cmdDB.AddCommand(&cobra.Command{
			Use:   "up-by-one",
			Short: "db up-by-one",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				if err := database.DB.UpByOne(ctx); err != nil {
					err = fmt.Errorf("db up-by-one error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			},
		})

		// db up-to command.
		cmdDB.AddCommand(&cobra.Command{
			Use:   "up-to",
			Short: "db up-to <version>",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				version, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					err = fmt.Errorf("db up-to error: parse arg %s: %w", args[0], err)
					errhandle.Handle(ctx, nil, err, true)
				}

				if err := database.DB.UpTo(ctx, version); err != nil {
					err = fmt.Errorf("db up-to error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			},
		})

		// db down command.
		cmdDB.AddCommand(&cobra.Command{
			Use:   "down",
			Short: "db down",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				if err := database.DB.Down(ctx); err != nil {
					err = fmt.Errorf("db down error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			},
		})

		// db down-to command.
		cmdDB.AddCommand(&cobra.Command{
			Use:   "down-to",
			Short: "db down-to <version>",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				version, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					err = fmt.Errorf("db down-to error: parse arg %s: %w", args[0], err)
					errhandle.Handle(ctx, nil, err, true)
				}

				if err := database.DB.DownTo(ctx, version); err != nil {
					err = fmt.Errorf("db down-to error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			},
		})

		// db status command.
		cmdDB.AddCommand(&cobra.Command{
			Use:   "status",
			Short: "db status",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				if err := database.DB.Status(); err != nil {
					err = fmt.Errorf("db status error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			},
		})

		// db create command.
		cmdDB.AddCommand(&cobra.Command{
			Use:   "create",
			Short: "db create <name>",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				if err := database.DB.Create(args[0]); err != nil {
					err = fmt.Errorf("db create error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			},
		})

		// db reset command.
		cmdDB.AddCommand(&cobra.Command{
			Use:   "reset",
			Short: "db reset",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				if err := database.DB.Reset(ctx); err != nil {
					err = fmt.Errorf("db reset error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			},
		})

		// Goose seed command.
		cmdDB.AddCommand(&cobra.Command{
			Use:   "seed",
			Short: "db seed",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				if err := database.DB.Seed(ctx); err != nil {
					err = fmt.Errorf("db reset error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			},
		})
	}

	// Server commands.
	// Start the server. Run function is not required
	// since it just needs to bypass the other commands.
	cmdRoot.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "start the server",
		Run: func(cmd *cobra.Command, args []string) {
			// Load env variables.
			env.Init()

			// Initialize aws client.
			if err := aws.Init(); err != nil {
				err = fmt.Errorf("aws init error: %w", err)
				errhandle.Handle(ctx, nil, err, true)
			}

			// Load config.
			config.Load()

			// Set up logger.
			logWriter, err := logger.New(logger.Config{
				Level:  config.Log.Level,
				Mode:   config.Log.Mode,
				SHASUM: SHASUM,

				Syslog: config.Log.Syslog,
				File:   config.Log.File,
				AWS:    config.Log.AWS,
			})
			if err != nil {
				err = fmt.Errorf("logger new error: %w", err)
				errhandle.Handle(ctx, nil, err, true)
			}

			// Set the logger.
			logger.Logger = logWriter

			// Initialize database.
			if database.DB, err = database.Connect(logger.Logger, config.Database); err != nil {
				err = fmt.Errorf("database new error: %w", err)
				errhandle.Handle(ctx, nil, err, true)
			}

			// Run migrations if the automigrate option is set.
			if config.Database.Migrations.Auto {
				if err := database.DB.Up(ctx); err != nil {
					err = fmt.Errorf("migrations auto error: db up error: %w", err)
					errhandle.Handle(ctx, nil, err, true)
				}
			}

			// Create the handler.
			handler := handler.New(handler.Config{
				SHASUM: SHASUM,
			})

			// Create the fiber app.
			app := fiber.New(fiber.Config{
				DisableStartupMessage: !env.IsDev(),
				ReadTimeout:           duration.Seconds(config.Server.Timeout.Read),
				WriteTimeout:          duration.Seconds(config.Server.Timeout.Write),
				IdleTimeout:           duration.Seconds(config.Server.Timeout.Idle),
				ErrorHandler:          handler.Error,
			})

			// Set up the middleware and router.
			middleware.Setup(app, handler)
			router.Setup(app, handler)

			// Listen.
			if err := app.Listen(config.Server.Addr); err != nil {
				err = fmt.Errorf("server start error: %w", err)
				errhandle.Handle(ctx, nil, err, true)
			}
			defer app.Shutdown()
		},
	})

	// Execute the root command.
	cmdRoot.Execute()
}
