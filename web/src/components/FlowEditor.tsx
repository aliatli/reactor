import React from 'react';
import ReactFlow, { 
    Controls, 
    Background,
    Connection 
} from 'reactflow';

interface FlowEditorProps {
    states: StateDefinition[];
    onSave: (flow: any) => void;
}

export const FlowEditor: React.FC<FlowEditorProps> = ({ states, onSave }) => {
    const [nodes, setNodes] = useState([]);
    const [edges, setEdges] = useState([]);

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