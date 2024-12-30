import React, { useState, useCallback } from 'react';
import ReactFlow, { 
    Controls, 
    Background,
    Connection,
    addEdge,
    Edge,
    Node,
    OnNodesChange,
    OnEdgesChange,
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

    return (
        <div style={{ height: '100vh' }}>
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