Objetivo: Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.

Clone o seguinte repositório: clique para acessar o repositório.

Toda rotina de criação do leilão e lances já está desenvolvida, entretanto, o projeto clonado necessita de melhoria: adicionar a rotina de fechamento automático a partir de um tempo.

Para essa tarefa, você utilizará o go routines e deverá se concentrar no processo de criação de leilão (auction). A validação do leilão (auction) estar fechado ou aberto na rotina de novos lançes (bid) já está implementado.

Você deverá desenvolver:

Uma função que irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente;
Uma nova go routine que validará a existência de um leilão (auction) vencido (que o tempo já se esgotou) e que deverá realizar o update, fechando o leilão (auction);
Um teste para validar se o fechamento está acontecendo de forma automatizada;

Dicas:

Concentre-se na no arquivo internal/infra/database/auction/create_auction.go, você deverá implementar a solução nesse arquivo;
Lembre-se que estamos trabalhando com concorrência, implemente uma solução que solucione isso:
Verifique como o cálculo de intervalo para checar se o leilão (auction) ainda é válido está sendo realizado na rotina de criação de bid;
Para mais informações de como funciona uma goroutine, clique aqui e acesse nosso módulo de Multithreading no curso Go Expert;
 
Entrega:

O código-fonte completo da implementação.
Documentação explicando como rodar o projeto em ambiente dev.
Utilize docker/docker-compose para podermos realizar os testes de sua aplicação.

### RODANDO EM AMBIENTE DEV

Para rodar o projeto em ambiente dev, siga os passos abaixo:

1. Clone o repositório:

2. Acesse a pasta do projeto:

3. Existe um arquivo .env na pasta cmd/auction. 
As Seguintes variáveis podem ser alteradas para configurar o tempo de fechamento do leilão:
- AUCTION_DURATION=10s #tempo de duração do leilão 
- CHECK_INTERVAL=10s #tempo de verificação de leilões para encerrar

4. Execute o comando abaixo para subir o container:

```bash
docker-compose up -d --build
```

5. Na raiz do projeto existe uma pasta **http** com um arquivo **test.http** que contém as requisições para testar a aplicação.

A seguir um exemplo de requisição para criar um leilão:
```bash
POST http://localhost:8080/auction
Content-Type: application/json

{
    "product_name": "Laptop 2",
    "category": "Electronics 2",
    "description": "A laptop with 8GB RAM and 1TB storage 2",
    "condition": 1
}

```
_______

A seguir uma requisição para recuperar o produto com o ID do leilão criado:
```bash
GET http://localhost:8080/auction?status=0&category=Electronics&productName=Celular
Content-Type: application/json
```
_______

A seguir uma requisição para recuperar um leilão específico:
```bash
GET http://localhost:8080/auction/b66301db-db75-4129-96d4-6da3d6f6a354
Content-Type: application/json
```

_______

A seguir uma requisição para dar um lance em um leilão:
```bash
POST http://localhost:8080/bid
Content-Type: application/json

{
    "user_id": "8bbd6e8a-3718-47cf-a587-0f5660373ccb", #utilize o ID do usuário que deseja dar o lance
    "auction_id": "8bbd6e8a-3718-47cf-a587-0f5660373ccb",
    "amount": 100
}
```

### Executando o teste 

Para executar o teste, execute o comando abaixo:

```bash
docker exec -it auction go test ./...
```
No teste está mocado na variável de ambiente o tempo de 10s para a duração do leilão e 10s para o intervalo de verificação de leilões para encerrar.