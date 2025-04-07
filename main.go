package main

import (
	"log"

	_ "github.com/lib/pq" // Add this import for PostgreSQL driver
)

func main() {
	redis := Redis{
		Host:     "localhost",
		Port:     "6379",
		Password: "", // Set if you have a password
		DB:       0,
	}

	pg := PostgreSQL{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres_database_testing_poc",
		Password: "postgres_database_testing_poc",
		DB:       "postgres_database_testing_poc",
	}

	userService := NewUserService(&pg, &redis)

	ginServer := GinServer{
		Host: "localhost",
		Port: "8080",
	}

	// Connect to database
	db, err := pg.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pg.Close(db)

	// Create table and verify it exists
	if err := pg.CreateUsersTable(db); err != nil {
		log.Fatal("Failed to create table:", err)
	}
	log.Println("Database table created/verified successfully")

	// Test Redis connection
	if redisClient, err := redis.Connect(); err == nil {
		defer redis.Close(redisClient)
		log.Println("Redis connection successful")
	} else {
		log.Printf("Warning: Redis connection failed: %v", err)
	}

	// Start the server
	log.Printf("Starting server on %s:%s", ginServer.Host, ginServer.Port)
	ginServer.Start(userService)
}
