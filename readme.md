# Rate Limiter API Example

Este projeto é uma API de exemplo para limitação de taxa (Rate Limiter) usando Redis. A limitação de taxa é uma técnica para limitar o número de solicitações que um cliente pode fazer para uma API em um determinado período de tempo.

## Como rodar a aplicação com Docker Compose

1. Certifique-se de que o Docker e o Docker Compose estão instalados em sua máquina.
2. Clone este repositório.
3. Navegue até o diretório do projeto.
4. Execute o comando `docker-compose up`.

O Docker Compose irá construir a imagem do Docker para a aplicação e iniciar os containers para a aplicação e o Redis.

## Funcionamento do Rate Limiter

O Rate Limiter pode limitar as solicitações por IP ou por token. As configurações para a limitação de taxa podem ser ajustadas no arquivo `.env`.

- **LIMIT_REQUESTS_DEFAULT_BY_IP**: O número máximo de solicitações que um IP pode fazer em um período de tempo especificado.
- **REQUEST_LIMIT_IN_SEC**: O período de tempo (em segundos) para o limite de solicitações por IP ou Token.
- **BLOCK_DURATION**: A duração (em segundos) que um IP ou Token será bloqueado após exceder o limite de solicitações.
- **LIMIT_REQUESTS_BY_TOKEN**: O número máximo de solicitações que um token pode fazer em um período de tempo especificado.
- **EXPIRATION_TOKEN**: A duração (em segundos) que um token é válido.

## Swagger

A documentação da API está disponível no Swagger. Após iniciar a aplicação, você pode acessar a documentação do Swagger em [http://localhost:8080/swagger-ui/index.html](http://localhost:8080/swagger-ui/index.html).

## Redis Insight

Para uma melhor visualização e gerenciamento do Redis, você pode instalar o Redis Insight. Siga as instruções de instalação na [página oficial do Redis Insight](https://redis.com/redis-enterprise/redis-insight/).

## Configuração

As configurações da aplicação podem ser ajustadas no arquivo `.env`. Este arquivo contém várias variáveis de ambiente que são usadas para configurar o comportamento da aplicação. As variáveis de ambiente incluem as configurações para a limitação de taxa, bem como a chave secreta usada para a autenticação.
