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
    Node,
    Edge as ReactFlowEdge
} from 'reactflow';
import 'reactflow/dist/style.css';
import { StateDefinition } from '../types/flow';
import { PrimitivePanel } from './PrimitivePanel';
import { Edge as CustomEdge } from '../types/flow';

interface FlowEditorProps {
    states: StateDefinition[];
    onSave: (flow: any) => void;
    onStateCreated?: () => void;
}

// Custom node component
const StateNode: React.FC<NodeProps> = ({ data, id }) => (
    <div 
        style={{ padding: '10px', border: '1px solid #ccc', borderRadius: '5px', background: 'white' }}
        onClick={(e) => {
            e.stopPropagation();
            data.onSelect(id);
        }}
    >
        <button 
            onClick={(e) => {
                e.stopPropagation();
                data.onDelete(id);
            }}
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
    const [edges, setEdges, onEdgesChange] = useEdgesState<CustomEdge[]>([]);
    const [showNewStateForm, setShowNewStateForm] = useState(false);
    const [newStateName, setNewStateName] = useState('');
    const [selectedState, setSelectedState] = useState<string | null>(null);
    const [primitives, setPrimitives] = useState<string[]>([]);
    const [isInitialized, setIsInitialized] = useState(false);

    const handleStateSelect = (stateId: string) => {
        setSelectedState(stateId);
    };

    // Memoize nodeTypes
    const nodeTypes = useMemo(() => ({
        stateNode: StateNode,
    }), []);

    const handleDeleteState = useCallback((stateId: string) => {
        // Delete from backend
        fetch(`http://localhost:8080/api/states/${stateId}`, {
            method: 'DELETE',
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to delete state');
            }

            // First update all connected states
            const updatePromises = states.map(state => {
                if (state.name !== stateId && (
                    state.transitions.success === stateId || 
                    state.transitions.failure === stateId ||
                    state.edges?.some(e => e.source === stateId || e.target === stateId)
                )) {
                    const updatedState: StateDefinition = {
                        ...state,
                        edges: (state.edges || []).filter(edge => 
                            edge.source !== stateId && edge.target !== stateId
                        ),
                        transitions: {
                            success: state.transitions.success === stateId ? 'none' : state.transitions.success,
                            failure: state.transitions.failure === stateId ? 'none' : state.transitions.failure
                        }
                    };

                    return fetch('http://localhost:8080/api/states', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(updatedState),
                    });
                }
                return Promise.resolve();
            });

            // Wait for all updates to complete
            return Promise.all(updatePromises);
        })
        .then(() => {
            // Then update the UI
            setNodes(nodes => nodes.filter(node => node.id !== stateId));
            setEdges(edges => edges.filter(edge => 
                edge.source !== stateId && edge.target !== stateId
            ));
            // Trigger a refresh of the states
            onStateCreated?.();
        })
        .catch(error => {
            console.error('Error deleting state:', error);
            alert('Failed to delete state');
        });
    }, [setNodes, setEdges, states, onStateCreated]);

    // Memoize the state conversion
    const stateNodes = useMemo(() => 
        states.map((state) => ({
            id: state.name,
            type: 'stateNode',
            data: { 
                label: state.name,
                onDelete: handleDeleteState,
                onSelect: handleStateSelect
            },
            position: state.position,
            draggable: true
        })), [states, handleDeleteState, handleStateSelect]);

    const stateEdges = useMemo(() => 
        states.flatMap(state => 
            (state.edges || [])
                .map(edge => ({
                    id: `${edge.source}-${edge.target}`,
                    source: edge.source,
                    target: edge.target,
                    sourceHandle: edge.sourceHandle,
                    style: { stroke: edge.sourceHandle === 'success' ? '#4CAF50' : '#f44336' },
                    label: edge.sourceHandle === 'success' ? 'Success' : 'Failure'
                }))
        ), [states]);

    // Update nodes and edges when states change
    useEffect(() => {
        if (!isInitialized && states.length > 0) {
            console.log('Initializing flow with states:', states);
            const initialNodes = states.map((state) => ({
                id: state.name,
                type: 'stateNode',
                data: { 
                    label: state.name,
                    onDelete: handleDeleteState,
                    onSelect: handleStateSelect
                },
                position: state.position,
                draggable: true
            }));

            // Collect all edges from all states
            const initialEdges = states.flatMap(state => 
                (state.edges || []).map(edge => ({
                    id: `${edge.source}-${edge.target}`,
                    source: edge.source,
                    target: edge.target,
                    sourceHandle: edge.sourceHandle,
                    type: 'default',
                    animated: false,
                    style: { stroke: edge.sourceHandle === 'success' ? '#4CAF50' : '#f44336' },
                    label: edge.sourceHandle === 'success' ? 'Success' : 'Failure'
                }))
            );

            console.log('Setting initial edges:', initialEdges);
            setNodes(initialNodes);
            setEdges(initialEdges);
            setIsInitialized(true);
        }
    }, [states, handleDeleteState, handleStateSelect, setNodes, setEdges, isInitialized]);

    // Add this effect to handle updates to states after initialization
    useEffect(() => {
        if (isInitialized && states.length > 0) {
            const currentEdges = states.flatMap(state => 
                (state.edges || []).map(edge => ({
                    id: `${edge.source}-${edge.target}`,
                    source: edge.source,
                    target: edge.target,
                    sourceHandle: edge.sourceHandle,
                    type: 'default',
                    animated: false,
                    style: { stroke: edge.sourceHandle === 'success' ? '#4CAF50' : '#f44336' },
                    label: edge.sourceHandle === 'success' ? 'Success' : 'Failure'
                }))
            );
            setEdges(currentEdges);
        }
    }, [states, isInitialized, setEdges]);

    // Add a cleanup effect
    useEffect(() => {
        return () => {
            setNodes([]);
            setEdges([]);
            setIsInitialized(false);
        };
    }, [setNodes, setEdges]);

    useEffect(() => {
        // Fetch primitives from backend
        fetch('http://localhost:8080/api/primitives')
            .then(res => res.json())
            .then(data => setPrimitives(data));
    }, []);

    const handlePrimitiveSave = (selectedPrimitives: string[]) => {
        if (!selectedState) return;

        const state = states.find(s => s.name === selectedState);
        if (!state) return;

        const updatedState: StateDefinition = {
            ...state,
            preliminaryActions: [{
                primitives: selectedPrimitives,
                executionOrder: 1
            }]
        };

        fetch('http://localhost:8080/api/states', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(updatedState),
        })
        .then(() => {
            onStateCreated?.();
            setSelectedState(null);
        })
        .catch(error => {
            console.error('Error saving primitives:', error);
            alert('Failed to save primitives');
        });
    };

    const onConnect = useCallback((connection: Connection) => {
        if (connection.source && connection.target && connection.sourceHandle) {
            const sourceNode = nodes.find(n => n.id === connection.source);
            if (!sourceNode) return;

            const isSuccess = connection.sourceHandle === 'success';
            const edge: CustomEdge = {
                id: `${connection.source}-${connection.target}`,
                source: connection.source,
                target: connection.target,
                sourceHandle: connection.sourceHandle,
                type: 'default',
                animated: false,
                style: { stroke: isSuccess ? '#4CAF50' : '#f44336' },
                label: isSuccess ? 'Success' : 'Failure'
            };

            const existingState = states.find(s => s.name === connection.source);
            if (!existingState) return;

            const stateDefinition: StateDefinition = {
                ...existingState,
                edges: [...(existingState.edges || []), {
                    id: `${connection.source}-${connection.target}`,
                    source: connection.source,
                    target: connection.target,
                    sourceHandle: connection.sourceHandle,
                    type: 'default',
                    animated: false,
                    style: { stroke: connection.sourceHandle === 'success' ? '#4CAF50' : '#f44336' },
                    label: connection.sourceHandle === 'success' ? 'Success' : 'Failure'
                }]
            };

            setEdges((eds) => [...eds, edge]);

            fetch('http://localhost:8080/api/states', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(stateDefinition),
            }).catch(error => {
                console.error('Error saving edge:', error);
            });
        }
    }, [nodes, states, setEdges]);

    const handleSave = () => {
        const stateDefinitions = nodes.reduce((acc, node) => {
            const nodeEdges = edges.filter(e => e.source === node.id).map(edge => ({
                id: edge.id,
                source: edge.source,
                target: edge.target,
                sourceHandle: edge.sourceHandle,
                type: edge.type,
                animated: edge.animated,
                style: edge.style,
                label: edge.label
            }));

            acc[node.id] = {
                name: node.id,
                position: {
                    x: node.position.x,
                    y: node.position.y
                },
                preliminaryActions: [],
                edges: nodeEdges,
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

        const position = {
            x: Math.random() * 500,
            y: Math.random() * 500
        };

        const newState: StateDefinition = {
            name: newStateName,
            preliminaryActions: [],
            position: position,
            edges: [],
            transitions: {
                success: 'none',
                failure: 'none'
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
                type: 'stateNode',
                data: { 
                    label: newStateName,
                    onDelete: handleDeleteState,
                    onSelect: handleStateSelect
                },
                position: position
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

    const onNodeDragStop = useCallback((_event: any, node: Node) => {
        // Find existing state to preserve its data
        const existingState = states.find(s => s.name === node.id);
        
        // Save the state with new position while preserving other data
        const stateDefinition: StateDefinition = {
            name: node.id,
            position: {
                x: node.position.x,
                y: node.position.y
            },
            preliminaryActions: existingState?.preliminaryActions || [],
            edges: existingState?.edges || [],
            transitions: {
                success: edges.find(e => e.source === node.id && e.sourceHandle === 'success')?.target || 'none',
                failure: edges.find(e => e.source === node.id && e.sourceHandle === 'failure')?.target || 'none'
            }
        };

        fetch('http://localhost:8080/api/states', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(stateDefinition),
        })
        .then(() => {
            onStateCreated?.();
        })
        .catch(error => {
            console.error('Error updating state position:', error);
        });
    }, [edges, states, onStateCreated]);

    return (
        <div style={{ width: '100%', height: '100%', position: 'relative' }}>
            <div style={{ 
                padding: '10px', 
                borderBottom: '1px solid #ccc',
                position: 'absolute',
                top: 0,
                left: 0,
                right: 0,
                zIndex: 4
            }}>
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
            
            <div style={{ 
                position: 'absolute',
                top: '60px',
                left: 0,
                right: 0,
                bottom: 0
            }}>
                <ReactFlow
                    nodes={nodes}
                    edges={edges}
                    onConnect={onConnect}
                    onNodesChange={onNodesChange}
                    onEdgesChange={onEdgesChange}
                    onNodeDragStop={onNodeDragStop}
                    nodeTypes={nodeTypes}
                    fitView
                    defaultViewport={{ x: 0, y: 0, zoom: 1 }}
                    minZoom={0.1}
                    maxZoom={4}
                    deleteKeyCode={['Backspace', 'Delete']}
                >
                    <Background />
                    <Controls />
                </ReactFlow>
            </div>

            {selectedState && (
                <PrimitivePanel
                    stateName={selectedState}
                    primitives={primitives}
                    selectedPrimitives={
                        states.find(s => s.name === selectedState)
                            ?.preliminaryActions[0]?.primitives || []
                    }
                    onClose={() => setSelectedState(null)}
                    onSave={handlePrimitiveSave}
                />
            )}
        </div>
    );
}; 