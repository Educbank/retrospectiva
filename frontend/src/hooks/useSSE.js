import { useEffect, useRef, useState } from 'react';
import toast from 'react-hot-toast';

const useSSE = (url, retrospectiveId) => {
  const [isConnected, setIsConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState(null);
  const eventSourceRef = useRef(null);

  useEffect(() => {
    if (!url || !retrospectiveId) return;

    const token = localStorage.getItem('token');
    if (!token) return;

    const sseUrl = `${url}?retrospective_id=${retrospectiveId}&token=${token}`;
    eventSourceRef.current = new EventSource(sseUrl);

    eventSourceRef.current.onopen = () => {
      console.log('SSE connected');
      setIsConnected(true);
    };

    eventSourceRef.current.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        setLastMessage(data);
        
        // Handle different message types
        switch (data.type) {
          case 'item_added':
            // Toast is handled by the mutation onSuccess
            break;
          case 'item_voted':
            // Toast is handled by the mutation onSuccess
            break;
          case 'item_deleted':
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
          case 'group_created':
            // Toast is handled by the mutation onSuccess
            break;
          case 'group_voted':
            // Toast is handled by the mutation onSuccess
            break;
          case 'group_deleted':
            // Toast is handled by the mutation onSuccess
            break;
          case 'items_merged':
            // Toast is handled by the mutation onSuccess
            break;
          case 'connected':
            console.log('Connected to retrospective:', data.data);
            break;
          case 'ping':
            // Keepalive, do nothing
            break;
          default:
            break;
        }
      } catch (error) {
        console.error('Error parsing SSE message:', error);
      }
    };

    eventSourceRef.current.onerror = (error) => {
      console.error('SSE error:', error);
      setIsConnected(false);
      
      // Attempt to reconnect after 3 seconds
      setTimeout(() => {
        if (eventSourceRef.current?.readyState === EventSource.CLOSED) {
          setIsConnected(true); // Trigger reconnection
        }
      }, 3000);
    };

    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
        setIsConnected(false);
      }
    };
  }, [url, retrospectiveId]);

  return {
    isConnected,
    lastMessage,
  };
};

export default useSSE;
