# Product Service API

Um microsserviço para gerenciamento de produtos implementado em Go, seguindo os princípios da Arquitetura Hexagonal e utilizando um banco de dados PostgreSQL para persistência.

## Arquitetura

Este projeto implementa a Arquitetura Hexagonal (também conhecida como Ports and Adapters). A estrutura do projeto é organizada em camadas bem definidas:

### Camadas da Arquitetura

1. **Core (Domínio)**: Contém as regras de negócio e entidades do domínio
   - `internal/core/models`: Entidades de domínio (Product)
   - `internal/core/ports`: Interfaces que definem os contratos entre as camadas

2. **Application**: Implementa os casos de uso da aplicação
   - `internal/application`: Serviços que orquestram as operações de negócio

3. **Adapters**: Implementações concretas das interfaces definidas no Core
   - **Driven Adapters** (saída): `internal/adapters/driven/postgresdb` - Implementação do repositório para PostgreSQL.
   - **Driver Adapters** (entrada): `internal/adapters/driver/http` - Handlers HTTP

### Benefícios desta Arquitetura

- **Testabilidade**: Facilita a criação de testes unitários isolados
- **Flexibilidade**: Permite trocar implementações (ex: banco de dados) sem afetar o core
- **Manutenibilidade**: Separação clara de responsabilidades
- **Independência de frameworks**: O domínio não depende de bibliotecas externas

## Estrutura do Projeto

```
product-service-go/
├── cmd/
│   └── main.go                 # Ponto de entrada da aplicação
├── internal/
│   ├── adapters/               # Camada de adaptadores
│   │   ├── driven/             # Adaptadores de saída (para infraestrutura)
│   │   │   ├── memdb/          # Implementação do repositório em memória (para testes)
│   │   │   └── postgresdb/     # Implementação do repositório com PostgreSQL
│   │   └── driver/
│   │       └── http/           # Handlers HTTP
│   ├── application/            # Serviços de aplicação
│   └── core/
│       ├── models/             # Entidades de domínio
│       └── ports/              # Interfaces (portas)
└── go.mod                      # Dependências do projeto
```

## API REST

O serviço expõe uma API REST para gerenciamento de produtos:

### Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/products` | Lista todos os produtos |
| GET | `/products/{id}` | Obtém um produto pelo ID |
| POST | `/products` | Cria um novo produto |
| PUT | `/products/{id}` | Atualiza um produto existente |
| DELETE | `/products/{id}` | Remove um produto |

### Formato dos Dados

**Produto (JSON)**:

```json
{
  "ID": "string",
  "Name": "string",
  "Price": 99.99,
  "CreatedAt": "string (ISO 8601)"
}
```

### Códigos de Status

- `200 OK`: Operação bem-sucedida
- `201 Created`: Recurso criado com sucesso
- `400 Bad Request`: Dados inválidos
- `204 No Content`: Operação bem-sucedida sem corpo de resposta
- `404 Not Found`: Recurso não encontrado
- `405 Method Not Allowed`: Método HTTP não suportado
- `500 Internal Server Error`: Erro interno do servidor

## Instalação e Execução

### Pré-requisitos

- Go 1.24 ou superior.
- PostgreSQL.

### Passos para Execução

1. Clone o repositório:

   ```bash
   git clone https://github.com/danielrios/product-service-go.git
   cd product-service-go
   ```

2. **Configure o Ambiente**:
   Crie um arquivo `.env` na raiz do projeto. Você pode copiar o exemplo:

   ```bash
   cp .env.example .env
   ```

   Edite o arquivo `.env` com as credenciais do seu banco de dados PostgreSQL, se forem diferentes do padrão.

3. **Prepare o Banco de Dados**:
   Conecte-se ao seu servidor PostgreSQL e execute os seguintes comandos para criar o banco de dados e a tabela:

   ```sql
   -- Cria o banco de dados
   CREATE DATABASE product_service_db;
   ```

   ```sql
   -- Conecte-se ao banco 'product_service_db' e crie a tabela
   CREATE TABLE products (
       id          TEXT PRIMARY KEY,
       name        TEXT NOT NULL,
       price       NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
       created_at  TIMESTAMPTZ NOT NULL
   );
   ```

4. Execute o serviço:

   ```bash
   go run cmd/main.go
   ```

5. O serviço estará disponível em `http://localhost:8080`

## Exemplos de Uso

### Listar todos os produtos

```bash
curl -X GET http://localhost:8080/products
```

### Obter um produto específico

```bash
curl -X GET http://localhost:8080/products/1
```

### Criar um novo produto

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"ID": "3", "Name": "Novo Produto", "Price": 299.99}'
```

### Atualizar um produto

```bash
curl -X PUT http://localhost:8080/products/3 \
  -H "Content-Type: application/json" \
  -d '{"ID": "3", "Name": "Produto Atualizado", "Price": 349.99}'
```

### Remover um produto

```bash
curl -X DELETE http://localhost:8080/products/3
```

## Desenvolvimento

### Executando Testes

```bash
go test -v ./...
```

## Características Técnicas

- **Persistência**: Utiliza **PostgreSQL** para armazenamento de dados, com o driver `pgx` de alta performance.
- **Roteamento HTTP**: Usa a biblioteca `chi` para um roteamento rápido, flexível e idiomático.
- **Configuração**: Carrega variáveis de ambiente a partir de um arquivo `.env` utilizando a biblioteca `godotenv`, facilitando o desenvolvimento local.
- **Graceful Shutdown**: Gerencia o encerramento adequado do servidor HTTP para não perder requisições em andamento, utilizando os pacotes `os/signal` e `context`.
- **Validação de Domínio**: Implementa validação de entidades diretamente no `core` da aplicação, garantindo a integridade dos dados.

## Contribuição

1. Fork o projeto
2. Crie sua branch de feature (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request
