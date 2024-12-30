import React, { useCallback, useEffect } from 'react';
import ReactFlow, { 
    Controls, 
    Background,
    Connection,
    addEdge,
    useNodesState,
    useEdgesState
} from 'reactflow';
import { StateDefinition } from '../types/flow';

interface FlowEditorProps {
    states: StateDefinition[];
    onSave: (flow: any) => void;
}

export const FlowEditor: React.FC<FlowEditorProps> = ({ states, onSave }) => {
    const [nodes, setNodes, onNodesChange] = useNodesState([]);
    const [edges, setEdges, onEdgesChange] = useEdgesState([]);

    const onConnect = useCallback((connection: Connection) => {
        setEdges((eds) => addEdge(connection, eds));
    }, []);

    useEffect(() => {
        // Convert states to nodes
        const stateNodes = states.map((state, index) => ({
            id: state.name,
            data: { label: state.name },
            position: { x: 100 * index, y: 100 }
        }));
        setNodes(stateNodes);
    }, [states]);

    return (
        <div style={{ height: '100vh' }}>
            <button onClick={() => onSave({ nodes, edges })}>Save Flow</button>
            <ReactFlow
                nodes={nodes}
                edges={edges}
                onConnect={onConnect}
                onNodesChange={onNodesChange}
                onEdgesChange={onEdgesChange}
            >
                <Background />
                <Controls />
            </ReactFlow>
        </div>
    );
}; 