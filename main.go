package main

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func main() {
	searchID := 5

	// Open Redis Connection
	redisClient := newRedisClient()
	result, err := redisPing(redisClient)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	// Search value in redis
	result, err = redisClient.Get("name_" + strconv.Itoa(searchID)).Result()
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT name FROM test WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(searchID)
	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("name is", name)
	}
}

func newRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return redisClient
}

func redisPing(client *redis.Client) (string, error) {
	result, err := client.Ping().Result()
	if err != nil {
		return "", err
	} else {
		return result, nil
	}
}
