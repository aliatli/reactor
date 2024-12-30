export interface Edge {
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