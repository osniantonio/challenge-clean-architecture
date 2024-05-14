# Go Expert

## Desafio Clean Architecture

### Executando

Antes de iniciar o programa, é necessário inciar o serviço de banco de dados. Para tanto, use o terminal para acessar a raiz deste projeto e então execute:

```sh
$ docker compose up -d
```

Se preferir, use o make:

```sh
$ make up
```

```sh
$ make down
```

O comando acima inicializará o banco de dados. Durante a inicialização, o diretório `docker-entrypoint-initdb.d` será copiado para o diretório homônimo dentro do servidor do banco de dados, garatindo a criação do banco de dados e da tabela para armazenamento dos dados referentes às odens de compra.

Para executar o programa, através do terminal acesse o diretório raiz deste projeto e execute os comandos a seguir:

```sh
$ cd cmd/ordersystem
$ go run main.go wire_gen.go
```

Se tudo correr bem, a seguinte mensagem será apresentada no terminal:

```
Starting web server on port :8000
Starting gRPC server on port 50051
Starting GraphQL server on port 8080
```

As portas usadas para inicializar os servidores podem ser alteradas através do arquivo `.env`, localizado no diretório `cmd/ordersystem`, relativo ao diretório raiz deste projeto.

### Inserindo novas ordens de compra

As ordens podem ser inseridas de três formas diferentes, explicadas a seguir.

#### Através de requisição HTTP

Para inserir uma nova ordem através de uma requisição HTTP, faça uma requisição POST para o endereço http://localhost:8000/order. O corpo da requisição deve conter um objeto JSON com o seguinte formato:

```json
{
  "id": "string",
  "price": "float",
  "tax": "float"
}
```

Também é possível realizar testes usando o arquivo `api.http`, localizado no diretório `api`, relativo ao diretório raiz deste projeto.

Se a requisição for bem sucedida, a reposta deverá retornar um objeto JSON com o seguinte formato:

```json
{
  "id": "string",
  "price": "float",
  "tax": "float",
  "final_price": "float"
}
```

### Através do GraphQL

Para criar uma nova ordem usando GraphQL, através de um navegador acesse o endereço http://localhost:8080. No _playground_ do GraphQL, insira o seguinte comando:

```graphql
mutation createOrder {
  createOrder(
    input: {
      id: "1822c225-f65b-44aa-afeb-c94b2a90ff78"
      price: 199.99
      tax: 11.8
    }
  ) {
    id
    price
    tax
    finalPrice
  }
}
```

Após executar a instrução, a resposta deverá ser um objeto JSON com o seguinte formato:

```json
{
  "data": {
    "createOrder": {
      "id": "1822c225-f65b-44aa-afeb-c94b2a90ff78",
      "price": 199.99,
      "tax": 11.8,
      "finalPrice": 211.79000000000002
    }
  }
}
```

### Através do gRPC

Para criar uma nova ordem usando gRPC, recomenda-se o uso do Evans (https://github.com/ktr0731/evans) como client gRPC. Para iniciar o Evans, através do terminal, digite:

```sh
$ evans -r repl
```

Após acessar o Evans, talvez seja necessário entrar no pacote `pb` e no serviço `OrderService`. Para tanto, execute os seguintes comandos:

```sh
> package pb
> service OrderService
```

Para efetuar a criação de uma nova ordem de compra, é necessário efetuar uma chamada RPC para `CreateOrder`. Para tanto, execute o seguinte comando:

```sh
call CreateOrder
```

O console do Evans solicitará informar os dados de `id`, `price` e `tax`. Informe os valores desejados e pressione ENTER. Se tudo correr bem, a resposta será um objeto JSON com os dados da ordem, como mostra o exemplo a seguir:

```sh
id (TYPE_STRING) => 605757f8-80d1-41cc-995b-08dd9a2aaac1 <ENTER>
price (TYPE_FLOAT) => 409.8 <ENTER>
tax (TYPE_FLOAT) => 9.8 <ENTER>
{
  "finalPrice": 419.59998,
  "id": "605757f8-80d1-41cc-995b-08dd9a2aaac1",
  "price": 409.8,
  "tax": 9.8
}
```

### Consultando ordens de compra

As ordens inseridas podem ser consultadas de três formas diferentes, explicadas a seguir.

#### Através de requisição HTTP

Para consultar todas as ordens de compra já inseridas, deve-se efetuar uma requisição GET para o endereço http://localhost:8000/order. Também é possível usar o arquivo `api.http`, localizado no diretório `api`, relativo ao diretório raiz deste projeto.

A resposta à requisição deverá ser um objeto JSON com todas as ordens inseridas no banco de dados.

```json
[
  {
    "Id": "a2d7cb34-2eb5-402a-879b-67940677c266",
    "Price": 599.09003,
    "Tax": 12.33,
    "FinalPrice": 611.42003
  },
  {
    "Id": "e767be01-d511-441d-90b2-38d9178779d1",
    "Price": 2999.99,
    "Tax": 39.8,
    "FinalPrice": 3039.79
  },
  {
    "Id": "fb86ce70-9fff-446a-a23d-fe2313500aca",
    "Price": 999.9,
    "Tax": 37.8,
    "FinalPrice": 1037.7
  }
]
```

### Através do GraphQL

Para consultar todas as ordens de compra através do GraphQL é necessário acessar o _playground_ do GraphQL usando o navegador, através do endereço http://localhost:8080. No campo de consulta, insira o seguinte comando:

```graphql
query ListOrders {
  orders {
    id
    price
    tax
    finalPrice
  }
}
```

Se tudo correr bem, a resposta deverá ser semelhante a esta:

```json
{
  "data": {
    "orders": [
      {
        "id": "a2d7cb34-2eb5-402a-879b-67940677c266",
        "price": 599.09003,
        "tax": 12.33,
        "finalPrice": 611.42003
      },
      {
        "id": "e767be01-d511-441d-90b2-38d9178779d1",
        "price": 2999.99,
        "tax": 39.8,
        "finalPrice": 3039.79
      },
      {
        "id": "fb86ce70-9fff-446a-a23d-fe2313500aca",
        "price": 999.9,
        "tax": 37.8,
        "finalPrice": 1037.7
      }
    ]
  }
}
```

### Através do gRPC

Para consultar as ordens de compra através de uma chamada gRPC, acesse o Evans. É importante garantir que o pacote e serviço estão corretamente selecionados. Portanto, execute os seguintes comandos:

```sh
> package db
> service OrderService
```

Para consultar todas as ordens disponívels, deve-se efetuar uma chamada RPC para `ListOrders`. Portanto, deve-se executar o seguinte:

```sh
> call ListOrders
```

Se tudo correr bem, a resposta será um objeto JSON contendo todas as ordens disponíveis, como mostrado no exemplo a seguir:

```sh
pb.OrderService@127.0.0.1:50052> call ListOrders
{
  "orders": [
    {
      "finalPrice": 611.42004,
      "id": "a2d7cb34-2eb5-402a-879b-67940677c266",
      "price": 599.09,
      "tax": 12.33
    },
    {
      "finalPrice": 3039.79,
      "id": "e767be01-d511-441d-90b2-38d9178779d1",
      "price": 2999.99,
      "tax": 39.8
    },
    {
      "finalPrice": 1037.7,
      "id": "fb86ce70-9fff-446a-a23d-fe2313500aca",
      "price": 999.9,
      "tax": 37.8
    }
  ]
}
```
