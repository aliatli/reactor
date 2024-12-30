export interface StateDefinition {
    name: string;
    preliminaryActions: PrimitiveChain[];
    mainAction?: string;
    transitions: {
        success: string;
        failure: string;
    };
}

export interface PrimitiveChain {
    primitives: string[];
    executionOrder: number;
} 