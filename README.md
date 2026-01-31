# Feedback Manager

Feedback Manager is a full-stack app built with **React (Vite)**, **Go**, and **PostgreSQL**, containerized with **Docker**.  
It provides CRUD APIs for employee feedback and a simple frontend to view, add, update, and delete entries.  
Environment variables manage configuration, with CORS enabled for smooth integration.

## Tech Stack
- Frontend: React + TypeScript + Vite
- Backend: Go + PostgreSQL
- Containerization: Docker + Docker Compose

## Setup
1. Clone the repo
2. Add `.env` files for backend and frontend
3. Run `docker-compose up --build`
4. Visit `http://localhost:5173` for frontend and `http://localhost:8080` for backend
