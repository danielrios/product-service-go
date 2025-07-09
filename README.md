# Product Service API

Um microsserviço para gerenciamento de produtos implementado em Go, seguindo os princípios da Arquitetura Limpa (Clean Architecture).

## Arquitetura

Este projeto implementa a Arquitetura Hexagonal (também conhecida como Ports and Adapters). A estrutura do projeto é organizada em camadas bem definidas:

### Camadas da Arquitetura

1. **Core (Domínio)**: Contém as regras de negócio e entidades do domínio
   - `internal/core/models`: Entidades de domínio (Product)
   - `internal/core/ports`: Interfaces que definem os contratos entre as camadas

2. **Application**: Implementa os casos de uso da aplicação
   - `internal/application`: Serviços que orquestram as operações de negócio

3. **Adapters**: Implementações concretas das interfaces definidas no Core
   - **Driven Adapters** (saída): `internal/adapters/driven/memdb` - Implementação do repositório
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
│   ├── adapters/
│   │   ├── driven/
│   │   │   └── memdb/          # Implementação do repositório em memória
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
- `404 Not Found`: Recurso não encontrado
- `405 Method Not Allowed`: Método HTTP não suportado
- `500 Internal Server Error`: Erro interno do servidor

## Instalação e Execução

### Pré-requisitos

- Go 1.24 ou superior

### Passos para Execução

1. Clone o repositório:
   ```bash
   git clone https://github.com/danielrios/product-service-go.git
   cd product-service-go
   ```

2. Execute o serviço:
   ```bash
   go run cmd/main.go
   ```

3. O serviço estará disponível em `http://localhost:8080`

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

- **Persistência**: Atualmente usa armazenamento em memória (pode ser substituído por um banco de dados real)
- **Concorrência**: Implementa mutex para operações thread-safe no repositório em memória
- **Graceful Shutdown**: Gerencia o encerramento adequado do servidor HTTP
- **Validação**: Implementa validação de entidades no domínio

## Contribuição

1. Fork o projeto
2. Crie sua branch de feature (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request
