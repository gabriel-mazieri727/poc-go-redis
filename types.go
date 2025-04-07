package main

import "github.com/gin-gonic/gin"

type PostgreSQL struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserService handles all user-related operations
type UserService struct {
	pg    *PostgreSQL
	redis *Redis
}

// Redis
type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// HTTP Server
type GinServer struct {
	Host   string
	Port   string
	Router *gin.Engine
}
