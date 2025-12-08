⚡ Go + Redis: Cache-Aside Pattern Demo

📺 Assista ao vídeo tutorial completo: [COLOQUE O LINK DO SEU VÍDEO AQUI]

Este projeto é uma demonstração prática de engenharia de performance. Ele mostra como escalar uma API que sofre com dependências lentas (600ms+) para uma resposta instantânea (<50ms) usando o padrão Cache-Aside com Go e Redis.

O projeto inclui uma UI de Benchmark ("Battle Mode") para visualizar a diferença de performance em tempo real, rodando diretamente no navegador.

🏗️ Arquitetura

O projeto segue o Standard Go Project Layout (Clean Architecture simplificada) para demonstrar como organizar projetos profissionais em Go.

/weather-cache
├── cmd/api/           # Entrypoint da aplicação (main.go)
├── internal/
│   ├── config/        # Gerenciamento de configuração (Env vars)
│   ├── handler/       # Camada HTTP (Request/Response & HTML ViewModels)
│   ├── service/       # Regra de Negócio (Lógica do Cache-Aside)
│   └── platform/      # Infraestrutura (Conexão Redis)
├── templates/         # Frontend (HTML + JS Benchmark)
└── Dockerfile         # Multi-stage build para Cloud Run


O Fluxo do Cache-Aside

sequenceDiagram
    participant User
    participant API as Go API
    participant Redis as Redis Cache
    participant Ext as External API (Lenta)

    User->>API: GET /weather
    API->>Redis: Tem essa chave?
    alt Cache HIT ⚡
        Redis-->>API: Sim! (Dados JSON)
        API-->>User: Retorna Imediato (< 5ms)
    else Cache MISS 🐢
        Redis-->>API: Não.
        API->>Ext: Busca dados (Lento ~500ms)
        Ext-->>API: Retorna dados
        API->>Redis: Salva dados (TTL 10s)
        API-->>User: Retorna dados
    end


🚀 Como Rodar Localmente

Pré-requisitos

Go 1.21+

Docker (para rodar o Redis localmente)

1. Subir o Redis

docker run --name my-redis -p 6379:6379 -d redis:alpine


2. Configurar e Rodar a API

# Clone o repositório
git clone [https://github.com/joaomarcosfurtado/go-redis-cache-aside-demo.git](https://github.com/joaomarcosfurtado/go-redis-cache-aside-demo.git)
cd go-redis-cache-aside-demo

# Baixe as dependências
go mod tidy

# Rode a aplicação
export REDIS_URL="redis://localhost:6379"
go run cmd/api/main.go


Acesse no navegador: http://localhost:8080

☁️ Deploy (Google Cloud Run)

Este projeto está pronto para Serverless. Usamos Google Cloud Run para a aplicação e Upstash (Redis Serverless).

Passo 1: Configurar Redis

Crie um banco no Upstash e copie a URL de conexão (redis://...).

Passo 2: Deploy

gcloud run deploy weather-app \
  --source . \
  --platform managed \
  --region us-east1 \
  --allow-unauthenticated \
  --set-env-vars REDIS_URL="SUA_URL_UPSTASH",REDIS_TLS="true"


📊 Teste de Carga (Benchmark)

Para provar a eficiência, o projeto possui um Modo de Simulação (?mock=true) que força um delay de 500ms no backend para simular uma API externa lenta sem ser bloqueado por Rate Limit.

Resultados Reais (Cloud Run + Upstash)

Métrica

Sem Cache (Miss) 🐢

Com Cache (Hit) ⚡

Melhoria

Latência (p95)

~700ms

~160ms*

4.3x Mais Rápido

RPS (Vazão)

~90 req/s

~300 req/s

3.3x Mais Capacidade

> Nota: 160ms representa a latência física de rede Brasil -> EUA. O tempo de processamento interno caiu para < 5ms.

Como reproduzir com k6

Instale o k6 e crie um arquivo loadtest.js:

import http from 'k6/http';
import { check } from 'k6';

export const options = { 
    vus: 50, 
    duration: '10s',
    thresholds: { http_req_duration: ['p(95)<200'] }
};

export default function () {
  // Troque pela sua URL do Cloud Run ou Localhost
  const url = '[https://SUA-APP.run.app/weather?lat=52&lon=13&mock=true](https://SUA-APP.run.app/weather?lat=52&lon=13&mock=true)';
  
  const res = http.get(url, {
      headers: { 'Accept': 'application/json' }
  });
  
  check(res, { 
      'status is 200': (r) => r.status === 200,
      'is hit': (r) => r.headers['X-Cache'] === 'HIT' 
  });
}


Rode o teste:

k6 run loadtest.js


🛠️ Tecnologias Utilizadas

Go 1.23: Backend de alta performance.

Redis (go-redis/v9): Armazenamento em memória chave-valor.

Google Cloud Run: Plataforma de container serverless.

Docker: Containerização multi-stage (imagem final Alpine).

k6: Ferramenta de Load Testing.

📝 Licença

Distribuído sob a licença MIT. Sinta-se livre para usar este código para estudos e projetos pessoais.

<p align="center">
Feito com 💜 por <a href="https://github.com/joaomarcosfurtado">João Marcos</a>
</p>