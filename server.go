package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (s *GinServer) Start(us *UserService) {
	r := gin.Default()
	r.SetTrustedProxies([]string{
		"127.0.0.1",
	})

	// Create API group and add routes to it
	api := r.Group("/api")
	{
		api.GET("/", us.fetchUsers)
		api.POST("/", us.createUser)
	}

	r.Run(fmt.Sprintf("%s:%s", s.Host, s.Port))
}
