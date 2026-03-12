# Taxi-app
Taxi App is a full-stack pet project that simulates a taxi ordering service with real-time driver location updates.

Users can create ride orders, track driver movement on the map, and view order history.
The system simulates driver movement by generating random coordinates and streaming them to clients via WebSockets.

This project demonstrates a microservice-style architecture using message queues and real-time communication.

# Features
* User authentication (JWT)

* Order creation and history

* Admin panel for managing drivers

* Real-time driver location updates

* Ride simulation with generated GPS coordinates

* WebSocket communication for live updates

* Message queue for driver location events

# Tech Stack

Frontend

* JavaScript

* React

Backend

* Go

* Gin Web Framework

Database

* PostgreSQL

Messaging

* RabbitMQ

Infrastructure

* Docker

* Nginx

# Running locally

``` 
git clone https://github.com/yourusername/taxi-app
docker compose -f docker-compose-dev.yml up --build
```
Application will be available at:
  http://localhost:3000

Swagger API docs:
  http://localhost:8081/docs
