import { useEffect, useRef, useState } from 'react';
import toast from 'react-hot-toast';

const useWebSocket = (url, retrospectiveId) => {
  const [isConnected, setIsConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState(null);
  const wsRef = useRef(null);
  const reconnectTimeoutRef = useRef(null);

  useEffect(() => {
    if (!url || !retrospectiveId) return;

    const connect = () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) return;

        const wsUrl = `${url}?retrospective_id=${retrospectiveId}&token=${token}`;
        wsRef.current = new WebSocket(wsUrl);

        wsRef.current.onopen = () => {
          console.log('WebSocket connected');
          setIsConnected(true);
          if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current);
            reconnectTimeoutRef.current = null;
          }
        };

        wsRef.current.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data);
            setLastMessage(message);
            
            // Handle different message types
            switch (message.type) {
              case 'item_added':
                // Toast is handled by the mutation onSuccess
                break;
              case 'item_voted':
                // Toast is handled by the mutation onSuccess
                break;
              case 'action_item_added':
                // Toast is handled by the mutation onSuccess
                break;
              case 'action_item_updated':
                // Toast is handled by the mutation onSuccess
                break;
              case 'action_item_deleted':
                // Toast is handled by the mutation onSuccess
                break;
              default:
                break;
            }
          } catch (error) {
            console.error('Error parsing WebSocket message:', error);
          }
        };

        wsRef.current.onclose = (event) => {
          console.log('WebSocket disconnected:', event.code, event.reason);
          setIsConnected(false);
          
          // Attempt to reconnect after 3 seconds
          if (event.code !== 1000) { // Don't reconnect if it was a normal closure
            reconnectTimeoutRef.current = setTimeout(() => {
              connect();
            }, 3000);
          }
        };

        wsRef.current.onerror = (error) => {
          console.error('WebSocket error:', error);
          setIsConnected(false);
        };

      } catch (error) {
        console.error('Error creating WebSocket connection:', error);
        setIsConnected(false);
      }
    };

    connect();

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close(1000, 'Component unmounting');
      }
    };
  }, [url, retrospectiveId]);

  const sendMessage = (message) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    }
  };

  return {
    isConnected,
    lastMessage,
    sendMessage,
  };
};

export default useWebSocket;
