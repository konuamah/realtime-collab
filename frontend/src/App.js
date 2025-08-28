import React, { useState, useEffect, useRef } from "react";

function App() {
  const [text, setText] = useState("");
  const [queue, setQueue] = useState([]); // unsent operations
  const ws = useRef(null);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    connectWebSocket();
  }, []);

  const connectWebSocket = () => {
    ws.current = new WebSocket("ws://localhost:8080/ws");

    ws.current.onopen = () => {
      setConnected(true);
      console.log("Connected to WebSocket");

      // Replay queued operations
      queue.forEach(op => ws.current.send(JSON.stringify(op)));
      setQueue([]); // clear queue
    };

    ws.current.onclose = () => {
      setConnected(false);
      console.log("Disconnected from WebSocket. Offline mode active.");
      // Will automatically reconnect
      setTimeout(connectWebSocket, 2000);
    };

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setText(data.text);
    };
  };


    const handleChange = (e) => {
    const newText = e.target.value;
    const cursorPos = e.target.selectionStart;
    let op;

    if (newText.length > text.length) {
      // insert
      op = { type: "insert", char: newText[cursorPos - 1], index: cursorPos - 1 };
    } else if (newText.length < text.length) {
      // delete
      op = { type: "delete", index: cursorPos };
    }

    setText(newText);

    if (connected && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(op));
    } else {
      // store operation for later
      setQueue(prev => [...prev, op]);
    }
  };


  return (
    <div style={{ padding: "20px" }}>
      <h2>Real-Time Collaboration Editor</h2>
      <textarea
        value={text}
        onChange={handleChange}
        rows={15}
        cols={80}
        style={{ fontSize: "16px" }}
      />
    </div>
  );
}

export default App;
