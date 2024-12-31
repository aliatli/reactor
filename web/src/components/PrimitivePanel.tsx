import React from 'react';

interface PrimitivePanelProps {
    stateName: string;
    primitives: string[];
    selectedPrimitives: string[];
    onClose: () => void;
    onSave: (primitives: string[]) => void;
}

export const PrimitivePanel: React.FC<PrimitivePanelProps> = ({
    stateName,
    primitives,
    selectedPrimitives,
    onClose,
    onSave,
}) => {
    const [selected, setSelected] = React.useState<string[]>(selectedPrimitives);

    return (
        <div style={{
            position: 'absolute',
            right: 0,
            top: 0,
            width: '300px',
            height: '100%',
            backgroundColor: 'white',
            boxShadow: '-2px 0 5px rgba(0,0,0,0.1)',
            padding: '20px',
            zIndex: 1000,
        }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '20px' }}>
                <h3>Primitives for {stateName}</h3>
                <button onClick={onClose}>×</button>
            </div>
            
            <div style={{ marginBottom: '20px' }}>
                {primitives.map(primitive => (
                    <div key={primitive} style={{ marginBottom: '10px' }}>
                        <label>
                            <input
                                type="checkbox"
                                checked={selected.includes(primitive)}
                                onChange={(e) => {
                                    if (e.target.checked) {
                                        setSelected([...selected, primitive]);
                                    } else {
                                        setSelected(selected.filter(p => p !== primitive));
                                    }
                                }}
                            />
                            {' '}{primitive}
                        </label>
                    </div>
                ))}
            </div>

            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px' }}>
                <button onClick={onClose}>Cancel</button>
                <button onClick={() => onSave(selected)}>Save</button>
            </div>
        </div>
    );
}; 