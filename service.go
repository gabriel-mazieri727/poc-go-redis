package main

func NewUserService(pg *PostgreSQL, redis *Redis) *UserService {
	return &UserService{
		pg:    pg,
		redis: redis,
	}
}
