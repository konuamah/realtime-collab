  # Real-Time Collaboration Engine (Google Docs Lite)

A lightweight full-stack application built with **Go** and **React** that enables collaborative document editing in real time, inspired by Google Docs. Multiple users can edit the same document simultaneously, with changes merged safely using **CRDTs**.

---

## Features

- Real-time collaboration across multiple users  
- Conflict-free state management using CRDT  
- WebSocket-based low-latency updates  
- Offline mode: edits are queued locally and synced automatically when connection is restored  
- Minimal React frontend for live editing  

---

## Tech Stack

- **Backend:** Go, Gorilla WebSocket, CRDT logic  
- **Frontend:** React, WebSocket API  
- **Concurrency & State Sync:** Thread-safe document updates using `sync.Mutex`  
- **Offline Support:** Local operation queue in React  

---

## Getting Started

### 1. Clone the repository
```
git clone https://github.com/konuamah/realtime-collab.git
cd realtime-collab
