Build da imagem do Docker:
docker build -t nomeImagem .

Rodar o container Redis:
docker run --name nomeContainer -p 6379:6379 -d nomeImagem

Verificar se o container está rodando:
docker ps



Criar o Ambiente Go:
go mod init nomeAmbiente
go get github.com/gin-gonic/gin
go get github.com/go-redis/redis/v8


Rodar o server go:
go run main.go



Testar:
for i in {1..15}; do curl -i http://localhost:8080/ping; done






docker exec -it containerId redis-cli
auth test123
keys *

