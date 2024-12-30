import React, { useEffect, useState } from 'react'
import { FlowEditor } from './components/FlowEditor'
import { StateDefinition } from './types/flow'

const App: React.FC = () => {
  const [states, setStates] = useState<StateDefinition[]>([])

  const fetchStates = () => {
    fetch('http://localhost:8080/api/states')
      .then(res => res.json())
      .then(data => setStates(Object.values(data)))
      .catch(error => console.error('Error fetching states:', error));
  }

  useEffect(() => {
    fetchStates()
  }, [])

  const handleSave = (flow: any) => {
    fetch('http://localhost:8080/api/flow', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ states: flow.nodes }),
    })
    .then(() => fetchStates())
    .catch(error => console.error('Error saving flow:', error));
  }

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <FlowEditor 
        states={states} 
        onSave={handleSave} 
        onStateCreated={fetchStates}
      />
    </div>
  )
}

export default App 