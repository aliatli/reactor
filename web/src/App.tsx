import React, { useEffect, useState, useCallback } from 'react'
import { FlowEditor } from './components/FlowEditor'
import { StateDefinition } from './types/flow'

function App() {
  const [states, setStates] = useState<StateDefinition[]>([])

  const fetchStates = useCallback(() => {
    fetch('http://localhost:8080/api/states')
      .then(res => res.json())
      .then((data: Record<string, StateDefinition>) => {
        console.log('Fetched states:', data);
        setStates(Object.values(data));
      })
      .catch(error => {
        console.error('Error fetching states:', error);
      });
  }, []);

  useEffect(() => {
    fetchStates();
  }, [fetchStates]);

  const handleSave = (flow: any) => {
    fetch('http://localhost:8080/api/flow', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(flow),
    })
    .then(() => {
      fetchStates();
    })
    .catch(error => {
      console.error('Error saving flow:', error);
    });
  };

  return (
    <div style={{ 
        width: '100vw', 
        height: '100vh', 
        position: 'relative', 
        overflow: 'hidden' 
    }}>
      <FlowEditor 
        states={states} 
        onSave={handleSave}
        onStateCreated={fetchStates}
      />
    </div>
  );
}

export default App; 