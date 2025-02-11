package main

import (
	"WIND/internal/database"
	"WIND/internal/handlers"
	"WIND/internal/taskService"
	"WIND/internal/userService"
	"WIND/internal/web/tasks"
	"WIND/internal/web/users"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	database.InitDB()
	if err := database.DB.AutoMigrate(&taskService.Task{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	taskRepo := taskService.NewTaskRepository(database.DB)
	tasksService := taskService.NewService(taskRepo)
	taskHandler := handlers.NewTaskHandler(tasksService)

	userRepo := userService.NewUserRepository(database.DB)
	usersService := userService.NewService(userRepo)
	userHandler := handlers.NewUserHandler(usersService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	taskStrictHandler := tasks.NewStrictHandler(taskHandler, nil)
	tasks.RegisterHandlers(e, taskStrictHandler)

	userStrictHandler := users.NewStrictHandler(userHandler, nil)
	users.RegisterHandlers(e, userStrictHandler)
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
