package main

import (
    "fmt"
	"net"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
)

type RedisClient struct {
    client *redis.Client
}

var (
    redisClient *RedisClient
)

const (
    redisAddress  = "localhost:6379"
    redisPassword = "test123"
)

func GetRedisClient() *RedisClient {
    redisClient = &RedisClient{
        client: redis.NewClient(&redis.Options{
            Addr:     redisAddress,
            Password: redisPassword,
            DB:       0,
        }),
    }
    return redisClient
}

func (rc *RedisClient) RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Capturar o IP
        IPAddress := c.GetHeader("X-Real-Ip")
        if IPAddress == "" {
            IPAddress = c.GetHeader("X-Forwarded-For")
        }
        if IPAddress == "" {
            IPAddress, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
			fmt.Println("Client IP Address:", IPAddress)
        }
        fmt.Println("Client IP Address:", IPAddress) // Log do IP

        // Carregar o script Lua
        script, err := os.ReadFile("script.lua")
        if err != nil {
            fmt.Println("Error reading script.lua:", err) // Log de erro
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "unable to read script"})
            return
        }

        // Preparar o script para execução
        takeScript := redis.NewScript(string(script))
        const rate = 10     // Tokens gerados por segundo
        const capacity = 10 // Capacidade máxima do bucket
        now := time.Now().UnixMicro() // Timestamp atual em microssegundos

        fmt.Println("Current Timestamp:", now) 

        // Executar o script Lua no Redis
        res, err := takeScript.Run(c, rc.client, []string{IPAddress}, capacity, rate, now, 1).Result()
		fmt.Println("Script result:", res)

        if err != nil {
            fmt.Println("Redis script execution error:", err)
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Redis script execution error"})
            return
        }

        // Inspecionar o resultado do script
        allowed, ok := res.(int64)
        if !ok {
            fmt.Println("Unexpected Redis response type:", res) 
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Unexpected Redis response type"})
            return
        }
        fmt.Println("Allowed (1 = yes, 0 = no):", allowed)

        // Bloquear requisições excedentes
        if allowed != 1 {
            fmt.Println("Request blocked for IP:", IPAddress) 
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"status": false, "message": "request overflowed"})
            return
        }

        // Permitir requisição
        fmt.Println("Request allowed for IP:", IPAddress) 
        c.Next()
    }
}

func PingHandler(c *gin.Context) {
    c.String(http.StatusOK, "Pong")
}

func main() {
    redisClient := GetRedisClient()
    router := gin.Default()
    router.Use(redisClient.RateLimitMiddleware())
    router.GET("/ping", PingHandler)

    if err := router.Run(":8080"); err != nil {
        fmt.Println("Failed to start server:", err)
    }
}
