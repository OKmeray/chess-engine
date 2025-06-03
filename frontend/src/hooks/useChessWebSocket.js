import { useState, useEffect } from "react";

const useChessWebSocket = (url, onMessageReceived) => {
    const [socket, setSocket] = useState(null);

    useEffect(() => {
        console.log("Connecting to WebSocket at:", url);

        const ws = new WebSocket(url);

        ws.onopen = () => console.log("WebSocket connection opened");
        ws.onmessage = (event) => {
            onMessageReceived(JSON.parse(event.data));
        };
        ws.onerror = (event) => console.error("WebSocket error:", event);
        ws.onclose = () => console.log("WebSocket closed");

        setSocket(ws);

        return () => {
            console.log("Cleaning up: closing WebSocket");
            ws.close();
        };
    }, [url]);

    const sendMessage = (data) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify(data));
        }
    };

    return { socket, sendMessage };
};

export default useChessWebSocket;
