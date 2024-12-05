import React, { useEffect, useState, useRef } from "react";

const WebSocketClient: React.FC = () => {
  const [serverTime, setServerTime] = useState("");
  const [connections, setConnections] = useState(0);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    socketRef.current = socket;

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setServerTime(data.server_time);
      setConnections(data.connections);
    };

    socket.onclose = () => {
      console.log("WebSocket disconnected");
    };

    return () => {
      socket.close();
    };
  }, []);

  return (
    <div style={{ padding: "20px", fontFamily: "Arial" }}>
      <h1>Server Dashboard</h1>
      <p>
        <strong>Server Time:</strong> {serverTime || "Connecting..."}
      </p>
      <p>
        <strong>Active Connections:</strong> {connections}
      </p>
    </div>
  );
};

export default WebSocketClient;
