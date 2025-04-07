package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler
func (s *UserService) fetchUsers(c *gin.Context) {
	start := time.Now()

	// First try to get from Redis
	redisClient, err := s.redis.Connect()
	if err != nil {
		log.Printf("Redis connection error: %v", err)
		// Continue with PostgreSQL if Redis fails
	} else {
		defer s.redis.Close(redisClient)

		// Try to get users from Redis
		usersJSON, err := s.redis.Get(redisClient, "users")
		if err == nil && usersJSON != "" {
			var users []User
			if err := json.Unmarshal([]byte(usersJSON), &users); err == nil {
				redisTime := time.Since(start)
				log.Printf("Redis cache hit - Time taken: %v", redisTime)
				c.JSON(200, gin.H{
					"cache": "hit",
					"time":  redisTime.String(),
					"users": users,
				})
				return
			}
		}
	}

	// If Redis fails or no data, fetch from PostgreSQL
	db, err := s.pg.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		c.JSON(500, gin.H{
			"message": "Error connecting to PostgreSQL",
			"error":   err.Error(),
		})
		return
	}
	defer s.pg.Close(db)

	var users []User
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Printf("Query error: %v", err)
		c.JSON(500, gin.H{
			"message": "Error fetching users",
			"error":   err.Error(),
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			c.JSON(500, gin.H{
				"message": "Error scanning user data",
				"error":   err.Error(),
			})
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		c.JSON(500, gin.H{
			"message": "Error processing results",
			"error":   err.Error(),
		})
		return
	}

	// Cache the results in Redis with a TTL of 5 minutes
	if redisClient, err := s.redis.Connect(); err == nil {
		defer s.redis.Close(redisClient)
		if usersJSON, err := json.Marshal(users); err == nil {
			// Set with 5 minute expiration
			s.redis.SetWithTTL(redisClient, "users", string(usersJSON), 5*time.Minute)
		}
	}

	totalTime := time.Since(start)
	log.Printf("PostgreSQL fetch - Time taken: %v", totalTime)
	c.JSON(200, gin.H{
		"users": users,
		"cache": "miss",
		"time":  totalTime.String(),
	})
}

func (s *UserService) createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Insert into PostgreSQL
	db, err := s.pg.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		c.JSON(500, gin.H{
			"message": "Error connecting to PostgreSQL",
			"error":   err.Error(),
		})
		return
	}
	defer s.pg.Close(db)

	_, err = db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
	if err != nil {
		log.Printf("Insert error: %v", err)
		c.JSON(500, gin.H{
			"message": "Error inserting user",
			"error":   err.Error(),
		})
		return
	}

	// Delete Redis cache instead of setting empty string
	if redisClient, err := s.redis.Connect(); err == nil {
		defer s.redis.Close(redisClient)
		s.redis.Delete(redisClient, "users")
	}

	c.JSON(200, gin.H{
		"message": "User created successfully",
	})
}
