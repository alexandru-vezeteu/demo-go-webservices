# POS Event Manager - Microservices Architecture

A robust event ticketing platform built with **Go microservices**, featuring Clean Architecture principles, gRPC/HTTP communication, and a React frontend.

> **Note on Database Design**: The database schema and entity relationships were imposed as part of the course requirements and do not represent the author's preferred design choices. This constraint leads to some awkward flows (e.g., ticket duplication across services, synchronous cross-service user creation) that would be structured differently in a greenfield project.

---

## Architecture Overview

The system is composed of **4 core microservices** communicating via gRPC and HTTP, with a React frontend:

| Service | Protocol | Port | Database | Description |
|---------|----------|------|----------|-------------|
| **IDM** | gRPC | 50051 | PostgreSQL + Redis | Identity Management - Authentication & JWT tokens |
| **IDM-Gateway** | HTTP | 8000 | - | REST-to-gRPC bridge for IDM |
| **User Service** | HTTP | - | MongoDB | User profiles & ticket management |
| **EventManager** | HTTP | - | PostgreSQL | Events, packets, tickets, inclusions |
| **Frontend** | HTTP | 5173 | - | React SPA |

### System Architecture Diagram

```
+-------------------------------------------------------------------------+
|                           pos-network                                    |
|  +------------------------------------------------------------------+   |
|  |  idm-network                                                      |   |
|  |  +-------------+     +-------------+     +-------------+        |   |
|  |  | IDM-Gateway |---->| IDM Service |---->| PostgreSQL  |        |   |
|  |  |  (HTTP:8000)| gRPC| (gRPC:50051)|     |   (IDM DB)  |        |   |
|  |  +-------------+     +------+------+     +-------------+        |   |
|  |                              |           +-------------+        |   |
|  |                              +---------->|    Redis    |        |   |
|  |                                          | (Blacklist) |        |   |
|  +------------------------------------------+-------------+--------+   |
|                                                                          |
|  +---------------------+    +---------------------+                    |
|  |  user-network       |    |  event-network       |                    |
|  |  +---------------+  |    |  +-----------------+|                    |
|  |  | User Service  |--+----+->|  EventManager   ||                    |
|  |  |    (HTTP)     |  |    |  |     (HTTP)      ||                    |
|  |  +-------+-------+  |    |  +--------+--------+|                    |
|  |  +-------v-------+  |    |  +--------v--------+|                    |
|  |  |    MongoDB    |  |    |  |   PostgreSQL    ||                    |
|  |  +---------------+  |    |  +-----------------+|                    |
|  +---------------------+    +---------------------+                    |
|                                                                          |
|  +-------------+                                                        |
|  |  Frontend   | (React:5173)                                           |
|  +-------------+                                                        |
+-------------------------------------------------------------------------+
```

---

## Design Choices and Advantages

### 1. Clean Architecture Pattern
All services follow Clean Architecture with clear layer separation:

| Layer | Responsibility | Benefits |
|-------|----------------|----------|
| **HTTP/gRPC Layer** | Request handling, routing | Protocol-agnostic business logic |
| **Application Layer** | Use cases, orchestration | Single responsibility per use case |
| **Domain Layer** | Business entities, interfaces | Zero external dependencies |
| **Infrastructure Layer** | DB implementations, external services | Easily swappable implementations |

**Why it matters:**
- **Testability** - Interfaces enable mocking for unit tests
- **Flexibility** - Swap databases without changing business logic
- **Maintainability** - Changes isolated to specific layers

### 2. gRPC for Internal Communication
IDM service uses gRPC instead of REST for:
- **Performance** - Binary protocol, faster serialization
- **Strong Contracts** - Proto definitions as source of truth
- **Type Safety** - Generated clients prevent runtime errors

### 3. JWT + Redis Token Blacklisting
- Stateless authentication with JWT tokens
- Redis stores revoked tokens (blacklist) with TTL matching token expiry
- **Fallback**: In-memory blacklist if Redis unavailable

### 4. Service-to-Service (S2S) Authentication
Dedicated service account (`serviciu_clienti` role) for inter-service calls:
- User Service -> EventManager authenticated requests
- Separate credentials from user tokens
- Isolated audit trail

### 5. Polyglot Persistence
Each service uses the optimal database:
| Service | Database | Rationale |
|---------|----------|-----------|
| IDM | PostgreSQL | ACID for credentials |
| EventManager | PostgreSQL | Relational data (events <-> packets) |
| User Service | MongoDB | Flexible document schema for profiles |

### 6. RBAC (Role-Based Access Control)
Four distinct roles with granular permissions:
- `admin` - System administrator
- `owner-event` - Event organizers (create events/packets)
- `client` - Regular customers (buy tickets)
- `serviciu_clienti` - Service account for S2S auth

---

## Known Limitations and Technical Debt

> **Important**: Several awkward patterns below (ticket duplication, synchronous cross-service calls during registration) stem from the imposed database schema requirements, not architectural choices.

| Issue | Impact | Recommendation |
|-------|--------|----------------|
| **DummyAuthz** in EventManager | Role checks not enforced | Implement proper RBAC middleware |
| **Ticket duplication** | Tickets stored in both MongoDB & PostgreSQL (imposed constraint) | Use EventManager as single source of truth |
| **No transaction handling** | Partial failures possible | Add distributed transaction support |
| **S2S sync calls** | Latency + failure cascading | Add circuit breaker + async fallback |
| **Missing audit logs** | No change tracking | Add structured logging |
| **No caching layer** | Repeated DB calls | Add Redis caching for hot data |
| **Cross-service user sync** | IDM synchronously calls User Service on registration (imposed flow) | Use event-driven async sync |

---

## What to Change / Add

### High Priority
1. **Implement proper RBAC middleware** - Replace `DummyAuthz` with real authorization checks
2. **Add circuit breaker** for inter-service calls (Hystrix/resilience pattern)
3. **Unified error handling** - Standardize error responses across services
4. **Add health check endpoints** - `/health` and `/ready` for Kubernetes

### Medium Priority
5. **Event-driven architecture** - Add Kafka/RabbitMQ for async operations
6. **Distributed tracing** - OpenTelemetry integration
7. **Caching layer** - Redis cache for frequently accessed data
8. **Rate limiting** - Protect APIs from abuse

### Future Enhancements
9. **API versioning** - `/v1/events` style versioning
10. **GraphQL Gateway** - Unified query interface for frontend
11. **Observability stack** - Prometheus + Grafana dashboards
12. **Saga pattern** - For distributed transactions (ticket purchase flow)

### If Database Design Were Not Imposed
- Single source of truth for tickets (EventManager only)
- Async user profile creation via message queue
- Denormalized read models for query optimization
- CQRS pattern for complex aggregates

---

## Architecture Diagrams

All diagrams are in PlantUML format in the `diagrams/` folder:

### Service Architecture
| Diagram | Description |
|---------|-------------|
| [idm_architecture.puml](diagrams/idm_architecture.puml) | IDM service Clean Architecture |
| [eventmanager_architecture.puml](diagrams/eventmanager_architecture.puml) | EventManager service layers |
| [user_service_architecture.puml](diagrams/user_service_architecture.puml) | User Service with adapters |

### Database and Infrastructure
| Diagram | Description |
|---------|-------------|
| [database_diagram.puml](diagrams/database_diagram.puml) | Full database schema across services |
| [docker_network_diagram.puml](diagrams/docker_network_diagram.puml) | Docker networking topology |

### Service Contracts (API)
| Diagram | Description |
|---------|-------------|
| [idm_service_contract.puml](diagrams/idm_service_contract.puml) | gRPC IDM endpoints |
| [eventmanager_service_contract.puml](diagrams/eventmanager_service_contract.puml) | HTTP EventManager API |
| [user_service_contract.puml](diagrams/user_service_contract.puml) | HTTP User Service API |

### User Flows (Sequence Diagrams)
| Diagram | Description |
|---------|-------------|
| [flow_register.puml](diagrams/flow_register.puml) | User registration flow |
| [flow_login.puml](diagrams/flow_login.puml) | Authentication flow |
| [flow_logout.puml](diagrams/flow_logout.puml) | Token revocation |
| [flow_buy_ticket.puml](diagrams/flow_buy_ticket.puml) | Ticket purchase with S2S auth |
| [flow_delete_event.puml](diagrams/flow_delete_event.puml) | Event deletion |
| [flow_delete_event_packet.puml](diagrams/flow_delete_event_packet.puml) | Packet deletion |
| [flow_delete_event_packet_inclusion.puml](diagrams/flow_delete_event_packet_inclusion.puml) | Inclusion removal |
| [flow_delete_ticket.puml](diagrams/flow_delete_ticket.puml) | Ticket deletion |
| [flow_update_event_seats.puml](diagrams/flow_update_event_seats.puml) | Seat capacity update |
| [flow_update_packet_seats.puml](diagrams/flow_update_packet_seats.puml) | Packet seats update |

### Security
| Diagram | Description |
|---------|-------------|
| [role_permissions.puml](diagrams/role_permissions.puml) | RBAC permissions matrix |

---

## Getting Started

### Prerequisites
- Docker
- Docker Compose

### Quick Start

```bash
# Clone the repository
git clone <repo-url>
cd demo-go-webservices

# Start all services
docker compose up -d

# View logs
docker compose logs -f
```

### Service URLs
| Service | URL |
|---------|-----|
| Frontend | http://localhost:5173 |
| IDM Gateway | http://localhost:8000/api/idm |
| User Service | http://localhost:8080/api/user-manager |
| EventManager | http://localhost:8081/api/event-manager |

---

## Project Structure

```
demo-go-webservices/
├── diagrams/              # PlantUML architecture diagrams
├── EventManager/          # Event management service (Go/Gin)
│   ├── app/
│   │   ├── domain/        # Business entities
│   │   ├── usecases/      # Application logic
│   │   ├── repositories/  # DB implementations
│   │   └── handlers/      # HTTP handlers
│   └── docker-compose.yaml
├── IDM/                   # Identity Management (Go/gRPC)
│   ├── app/
│   │   ├── domain/
│   │   ├── usecases/
│   │   ├── infrastructure/
│   │   └── server/        # gRPC server
│   └── docker-compose.yaml
├── IDM-Gateway/           # REST-to-gRPC bridge
├── User/                  # User profile service (Go/Gin)
│   ├── app/
│   │   ├── domain/
│   │   ├── usecases/
│   │   └── adapters/      # External service adapters
│   └── docker-compose.yaml
├── frontend/              # React SPA
└── docker-compose.yaml    # Root compose file
```

---

## API Quick Reference

### IDM (gRPC via Gateway)
```
POST /api/idm/auth/register  - Register new user
POST /api/idm/auth/login     - Login, get JWT token
POST /api/idm/auth/logout    - Revoke token
POST /api/idm/auth/verify    - Verify token validity
```

### EventManager
```
POST   /api/event-manager/events           - Create event
GET    /api/event-manager/events           - List/filter events
GET    /api/event-manager/events/:id       - Get event
PATCH  /api/event-manager/events/:id       - Update event
DELETE /api/event-manager/events/:id       - Delete event

(Similar CRUD for /event-packets, /tickets)
```

### User Service
```
POST   /api/user-manager/users             - Create user profile
GET    /api/user-manager/users/:id         - Get user
PATCH  /api/user-manager/users/:id         - Update user
DELETE /api/user-manager/users/:id         - Delete user

POST   /api/user-manager/clients/:id/tickets  - Buy ticket
```

---

## Security Model

### Authentication Flow
1. User registers -> IDM creates user in PostgreSQL + syncs to User Service (MongoDB)
2. User logs in -> IDM validates credentials, returns JWT
3. Requests include `Authorization: Bearer <token>` header
4. Services verify tokens via gRPC call to IDM

### Token Lifecycle
- JWT contains: `user_id`, `email`, `role`, `exp`
- On logout: Token added to Redis blacklist
- Verification checks both expiry AND blacklist

---

## Tech Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.21 |
| **Web Framework** | Gin |
| **RPC** | gRPC + Protocol Buffers |
| **ORM** | GORM (PostgreSQL), MongoDB Driver |
| **Auth** | JWT (golang-jwt) + bcrypt |
| **Cache** | Redis |
| **Frontend** | React + Vite |
| **Containers** | Docker + Docker Compose |

---

## License

This project is developed for academic purposes as part of the POS course.
