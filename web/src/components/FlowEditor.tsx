import React, { useCallback, useEffect, useState, useMemo } from 'react';
import ReactFlow, { 
    Controls, 
    Background,
    Connection,
    addEdge,
    useNodesState,
    useEdgesState,
    Handle,
    Position,
    NodeProps,
    Node
} from 'reactflow';
import 'reactflow/dist/style.css';
import { StateDefinition } from '../types/flow';

interface FlowEditorProps {
    states: StateDefinition[];
    onSave: (flow: any) => void;
    onStateCreated?: () => void;
}

// Custom node component
const StateNode: React.FC<NodeProps> = ({ data, id }) => (
    <div style={{ padding: '10px', border: '1px solid #ccc', borderRadius: '5px', background: 'white' }}>
        <button 
            onClick={() => data.onDelete(id)}
            style={{ 
                position: 'absolute', 
                right: -10, 
                top: -10,
                border: '1px solid #f44336',
                borderRadius: '50%',
                width: '20px',
                height: '20px',
                padding: 0,
                background: 'white',
                color: '#f44336',
                cursor: 'pointer'
            }}
        >
            ×
        </button>
        <Handle
            type="target"
            position={Position.Top}
            style={{ background: '#555' }}
        />
        <div>{data.label}</div>
        <Handle
            type="source"
            position={Position.Bottom}
            id="success"
            style={{ background: '#4CAF50', bottom: 10 }}
        />
        <Handle
            type="source"
            position={Position.Bottom}
            id="failure"
            style={{ background: '#f44336' }}
        />
    </div>
);

export const FlowEditor: React.FC<FlowEditorProps> = ({ states, onSave, onStateCreated }) => {
    const [nodes, setNodes, onNodesChange] = useNodesState([]);
    const [edges, setEdges, onEdgesChange] = useEdgesState([]);
    const [showNewStateForm, setShowNewStateForm] = useState(false);
    const [newStateName, setNewStateName] = useState('');

    // Memoize nodeTypes
    const nodeTypes = useMemo(() => ({
        stateNode: StateNode,
    }), []);  // Empty dependency array since StateNode never changes

    const handleDeleteState = useCallback((stateId: string) => {
        // Delete from backend
        fetch(`http://localhost:8080/api/states/${stateId}`, {
            method: 'DELETE',
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to delete state');
            }
            // Remove node
            setNodes(nodes => nodes.filter(node => node.id !== stateId));
            // Remove associated edges
            setEdges(edges => edges.filter(edge => 
                edge.source !== stateId && edge.target !== stateId
            ));
        })
        .catch(error => {
            console.error('Error deleting state:', error);
            alert('Failed to delete state');
        });
    }, [setNodes, setEdges]);

    useEffect(() => {
        const stateNodes = states.map((state) => ({
            id: state.name,
            type: 'stateNode',
            data: { 
                label: state.name,
                onDelete: handleDeleteState  // Pass delete handler to node
            },
            position: state.position || { x: Math.random() * 500, y: Math.random() * 500 }
        }));
        setNodes(stateNodes);
    }, [states, setNodes, handleDeleteState]);

    const onConnect = useCallback((connection: Connection) => {
        const sourceNode = nodes.find(n => n.id === connection.source);
        const targetNode = nodes.find(n => n.id === connection.target);
        
        if (sourceNode && targetNode) {
            const isSuccess = connection.sourceHandle === 'success';
            const edge = {
                id: `${connection.source}-${connection.target}`,
                source: connection.source,
                target: connection.target,
                sourceHandle: connection.sourceHandle,
                style: { stroke: isSuccess ? '#4CAF50' : '#f44336' },
                label: isSuccess ? 'Success' : 'Failure'
            };
            setEdges((eds) => [...eds, edge]);
        }
    }, [nodes, setEdges]);

    const handleSave = () => {
        const stateDefinitions = nodes.reduce((acc, node) => {
            acc[node.id] = {
                name: node.id,
                position: node.position,
                preliminaryActions: [],
                transitions: {
                    success: edges.find(e => e.source === node.id && e.sourceHandle === 'success')?.target || 'none',
                    failure: edges.find(e => e.source === node.id && e.sourceHandle === 'failure')?.target || 'none'
                }
            };
            return acc;
        }, {} as Record<string, any>);

        onSave({ states: stateDefinitions });
    };

    const handleAddState = () => {
        if (!newStateName.trim()) {
            alert('Please enter a state name');
            return;
        }

        const newState: StateDefinition = {
            name: newStateName,
            preliminaryActions: [],
            transitions: {
                success: '',
                failure: ''
            }
        };

        fetch('http://localhost:8080/api/states', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(newState),
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to create state');
            }
            return response.json();
        })
        .then(() => {
            // Add node to flow
            const newNode = {
                id: newStateName,
                data: { label: newStateName },
                position: { x: Math.random() * 500, y: Math.random() * 500 }
            };
            setNodes((nds) => [...nds, newNode]);
            setNewStateName('');
            setShowNewStateForm(false);
            onStateCreated?.();
        })
        .catch(error => {
            console.error('Error creating state:', error);
            alert('Failed to create state');
        });
    };

    return (
        <div style={{ width: '100%', height: '100%' }}>
            <div style={{ padding: '10px', borderBottom: '1px solid #ccc' }}>
                <button onClick={() => setShowNewStateForm(true)}>Add New State</button>
                <button onClick={handleSave} style={{ marginLeft: '10px' }}>Save Flow</button>
                
                {showNewStateForm && (
                    <div style={{ marginTop: '10px' }}>
                        <input
                            type="text"
                            value={newStateName}
                            onChange={(e) => setNewStateName(e.target.value)}
                            placeholder="Enter state name"
                        />
                        <button onClick={handleAddState} style={{ marginLeft: '5px' }}>Create</button>
                        <button onClick={() => setShowNewStateForm(false)} style={{ marginLeft: '5px' }}>Cancel</button>
                    </div>
                )}
            </div>
            
            <div style={{ height: 'calc(100% - 60px)' }}>
                <ReactFlow
                    nodes={nodes}
                    edges={edges}
                    onConnect={onConnect}
                    onNodesChange={onNodesChange}
                    onEdgesChange={onEdgesChange}
                    nodeTypes={nodeTypes}
                    fitView
                >
                    <Background />
                    <Controls />
                </ReactFlow>
            </div>
        </div>
    );
}; 