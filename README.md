# Rate Limiter com Redis e Golang

Este projeto implementa um **rate limiter** utilizando Redis e Golang. O objetivo é limitar o número de requisições que um cliente pode fazer em um determinado período de tempo, protegendo a API contra abusos.

## Passos para Configurar e Rodar o Projeto

```bash
# 1. Build da Imagem do Docker
# Crie a imagem do Redis utilizando o Dockerfile existente:
docker build -t nomeImagem .

# Rode o container Redis:
docker run --name nomeContainer -p 6379:6379 -d nomeImagem

# Verifique se o container está rodando:
docker ps

# 2. Criar o Ambiente Go
# Inicie o ambiente Go:
go mod init nomeAmbiente

# Instale as dependências do projeto:
# Instale o Gin:
go get github.com/gin-gonic/gin

# Instale o Redis:
go get github.com/go-redis/redis/v8

# 3. Rodar o Servidor Go
# Execute o servidor Golang:
go run main.go

# 4. Testar o Middleware
# Envie múltiplas requisições para testar o rate limiter:
for i in {1..15}; do curl -i http://localhost:8080/ping; done

# As primeiras 10 requisições retornarão "200 OK".
# As requisições excedentes retornarão "429 Too Many Requests".

# 5. Verificar o Redis
# Acesse o Redis para inspecionar as chaves armazenadas:
docker exec -it containerId redis-cli

# Faça login no Redis:
auth test123

# Liste as chaves armazenadas:
keys *

# Você verá chaves no formato:
# - 127.0.0.1:tokens: Representa a quantidade de tokens restantes para um cliente.
# - 127.0.0.1:last_access: Representa o timestamp do último acesso do cliente.


Resumo dos comandos

# Docker
docker build -t nomeImagem .
docker run --name nomeContainer -p 6379:6379 -d nomeImagem
docker ps

# Ambiente Go
go mod init nomeAmbiente
go get github.com/gin-gonic/gin
go get github.com/go-redis/redis/v8
go run main.go

# Testes
for i in {1..15}; do curl -i http://localhost:8080/ping; done

# Redis
docker exec -it containerId redis-cli
auth test123
keys *
