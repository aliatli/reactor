import { Edge as ReactFlowEdge } from 'reactflow';

export interface Edge extends Omit<ReactFlowEdge, 'sourceHandle'> {
    source: string;
    target: string;
    sourceHandle: string;
}

export interface StateDefinition {
    name: string;
    preliminaryActions: PrimitiveChain[];
    mainAction?: string;
    position: {
        x: number;
        y: number;
    };
    edges?: Edge[];
    transitions: {
        success: string;
        failure: string;
    };
}

export interface PrimitiveChain {
    primitives: string[];
    executionOrder: number;
} 