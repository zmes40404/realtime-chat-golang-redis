# 🧩 realtime-chat-golang-redis

A multi-room real-time chat server built with Golang, Gin, WebSocket, and Redis.

This side project is designed to explore modern backend engineering patterns including Pub/Sub messaging, WebSocket communication, and modular service architecture. It serves as a personal showcase for backend and DevOps skills.



## 🚀 Features

- ✅ Real-time messaging via WebSocket
- ✅ Support for multiple chat rooms
- ✅ Nickname-based user login
- ✅ In-memory message broadcasting
- 🚧 Redis Pub/Sub integration (coming soon)
- 🚧 Voting mechanism & photo upload (planned)



## 🧱 Tech Stack

- Language: **Golang**
- Web Framework: **Gin**
- Realtime: **WebSocket** (gorilla/websocket)
- Message Sync: **Redis Pub/Sub** (planned)
- State Mgmt: **Go Maps + Mutex**
- Deployment Target: **Docker + GitHub Actions** (future)


## 📁 Project Structure

realtime-chat-golang-redis/  
├── cmd/ main.go # Entry point  
├── internal/  
│ ├── server/ # HTTP + WebSocket server  
│ ├── chat/ # Hub, Client, Room logic  
│ └── redis/ # Redis connection (WIP)  
├── static/ # Frontend UI (HTML/JS)  
├── go.mod  
├── go.sum  
├── README.md


## 📌 Project Architecture Overview

```mermaid
graph TD
    A[main.go<br>App 入口點] --> B[handler.go<br>HTTP/WebSocket Router]
    B --> C{判斷請求<br>HTTP or WS}
    C -->|HTTP| D[回傳簡單 JSON 或頁面]
    C -->|WebSocket| E[建立 WS Conn<br>升級為 Client]
    E --> F[Client Struct<br>包含 conn / send chan]
    F --> G[註冊 Client 到<br>指定 Room]
    G --> H[Room Struct<br>管理 Client 列表 + 廣播 chan]
    H --> I[Hub Struct<br>管理所有 Room Map<br>+ Mutex 控制]
    H --> F1[廣播訊息給所有 Client]
    F1 --> F

    subgraph 外部通訊流程
        J[前端網頁 client.html<br>與 WebSocket 相連]
        J -->|連線請求| B
        F -->|傳入訊息| H
        H -->|廣播訊息| F1
        F1 -->|回送訊息| J
    end

    subgraph internal package
        F
        G
        H
        I
    end
```


## 🔧 Getting Started

### Prerequisites
- Go 1.20+ installed
- Redis (optional for Pub/Sub)

### Run locally

```bash
git clone https://github.com/zmes40404/realtime-chat-golang-redis.git
cd realtime-chat-golang-redis
go run ./cmd/main.go
```


## 📅 Development Plan

| Phase | Goal | Status |
| ----- | ---- | ------ |
| Phase 1 | Basic chatroom with WebSocket&multi-room | ✅In progress |
| Phase 2 | Add Redis Pub/Sub, vote, photo upload, admin roles | 🔜Coming soon |
| Phase 3 | Refactor and Code Optimization | 🚩 Future work


## ⌚ Future Work
| Item | Description | Status |
| ---- | ----------- | ------ |
|  1.  | Pull all the const setting to .env or config.yaml files | 🔜Coming soon |



## 🙋 Author
- Maton Wang (Pin Li Wang)  
- [LinkedIn](www.linkedin.com/in/matonwang)    
- Side project for Golang + Redis practice  
- For job applications in DevOps, backend, and Web3 sectors