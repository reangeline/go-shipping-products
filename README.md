# 📦 Shipping Products

This project solves the pack calculation problem.
Given an order with a quantity of items, the system returns the best combination of available packs, minimizing the total number of items and packages.

# 🔑 Features

- **API HTTP** in Go (clean architecture):
  - `GET /v1/packsizes` → lists the configured pack sizes.
  - `POST /v1/calculate` → calculates the optimal combination for an order.
- **Frontend React**:
  - Displays the available pack sizes.
  - Allows calculating packages for an order and visualizing the result.


## 🛠️ Technologies

- Go 1.24
- Gin (HTTP router)
- React
- Docker + Compose
- Nginx (frontend e proxy)
- Clean Architecture


## ⚙️ Requirements
- Docker 
- Make (for builds/tests convenience) 

### Variáveis de ambiente
  PACK_PROVIDER=file
  PACK_SIZES_FILE=./packs.csv
  HTTP_ADDR=:8080

## 🚀 How to Run

  Clone the repository:
   git clone https://github.com/reangeline/go-shipping-products.git 
   cd go-shipping-products
   git checkout main
   git pull

### Using Docker
  - Run all
    make all

#### Frontend: http://localhost:3000
#### Swagger Doc: http://localhost:8080/docs


### Locally (without Docker)
#### Backend
 - tests
  make api-test
 or
  go test ./... -v
____________________
  
  If your docker is running
   make docker-down

 - Run api
  make api-run
 or
  go run cmd/api/main.go

#### Obs: After run you have an api available in your http://localhost:8080
#### Swagger Doc: http://localhost:8080/docs

#### Frontend
  - Run Frontend
    make web-dev
  or
    cd web && npm install && npm run dev

#### Available to access in http://localhost:5173

# ✏️ Logic proccess to development
1.	I used Clean Architecture and started by creating the core with all the business logic.
1.1 First, I built the order.go entity, which handles the order quantity requested by the customer.
1.2 Then I created pack.go, which defines the fixed (immutable) pack sizes.
1.3 In both cases, I added a simple rule: the values must always be greater than zero and wrote tests for them.
1.4 After that, I worked on the most important part — the calculator.
1.4.1 The calculator takes the order quantity and the list of available packs, and applies the algorithm to find the best combination.
1.4.2 The algorithm works with the pack sizes and the target quantity, and is optimized using GCD.
1.4.3 While creating the calculator, I used Table-Driven Tests to validate the whole process.

2.	After the domain, I created ports/inbound/order for my use case contracts.
2.1 Next, I created outbound/packagesizes with a provider interface to list values. In this case, I could also add a repository interface to fetch values from a database.

3.	Then I implemented my use cases and finished my core application.

4.	For adapters, I followed the same logic with inbound and outbound. Nothing special — the main difference was abstracting the Gin framework so it can be replaced by others.
4.1 I created external communication using HTTP, but this could also be gRPC, CLI, GraphQL, etc.
4.2 In outbound, I created packsizes/file provider to read from a CSV file. I could also use an env provider (to read from environment variables) or a repository provider (to read from a database). This is the external dependency of my use cases.
4.3 I added documentation using Swagger (OpenAPI) for the app.

5.	Finally, I created the wiring to handle dependency injection.
6.	I built a simple frontend using React + Vite.
7.	And the Dockerfiles/docker-compose to run the whole project.


👨‍💻 Author
    Renato Angeline
    Senior Software Engineer




