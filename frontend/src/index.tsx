import { render } from 'preact';
import { useState } from 'preact/hooks';
import './style.css';

export function App() {
	const [numVariables, setNumVariables] = useState(2);
	const [numConstraints, setNumConstraints] = useState(2);
	const [objective, setObjective] = useState([1, 1]);
	const [objectiveType, setObjectiveType] = useState<'maximize' | 'minimize'>('maximize');
	const [constraints, setConstraints] = useState([[1, 1, 1], [1, 2, 2]]);
	const [constraintSigns, setConstraintSigns] = useState<string[]>(Array(2).fill("<="));
	const [result, setResult] = useState<any>(null);
	const [loading, setLoading] = useState(false);

	const updateObjective = (index: number, value: string) => {
		const newObjective = [...objective];
		newObjective[index] = parseFloat(value) || 0;
		setObjective(newObjective);
	};

	const updateConstraint = (row: number, col: number, value: string) => {
		const newConstraints = [...constraints];
		newConstraints[row][col] = parseFloat(value) || 0;
		setConstraints(newConstraints);
	};

	const adjustDimensions = (newVars, newConstraints) => {
		// Actualizar función objetivo
		const newObjective = Array(newVars).fill(0);
		for (let i = 0; i < Math.min(objective.length, newVars); i++) {
			newObjective[i] = objective[i];
		}
		setObjective(newObjective);

		// Actualizar restricciones
		const newConstraintsMatrix = Array(newConstraints).fill(null).map(() => Array(newVars + 1).fill(0));
		for (let i = 0; i < Math.min(constraints.length, newConstraints); i++) {
			for (let j = 0; j < Math.min(constraints[i].length, newVars + 1); j++) {
				newConstraintsMatrix[i][j] = constraints[i][j];
			}
		}
		setConstraints(newConstraintsMatrix);

		// Ajustar signos de restricciones cuando cambia el número de restricciones
		const newSigns = Array(newConstraints).fill("<=");
		for (let i = 0; i < Math.min(constraintSigns.length, newConstraints); i++) {
			newSigns[i] = constraintSigns[i];
		}
		setConstraintSigns(newSigns);
	};

	const handleVariablesChange = (value: string) => {
		const newVars = parseInt(value) || 2;
		setNumVariables(newVars);
		adjustDimensions(newVars, numConstraints);
	};

	const handleConstraintsChange = (value: string) => {
		const newConstraints = parseInt(value) || 2;
		setNumConstraints(newConstraints);
		adjustDimensions(numVariables, newConstraints);
	};

	const solveSimplex = async () => {
		setLoading(true);
		setResult(null);

		const requestBody = {
			objective: {
				n: numVariables,
				coefficients: objective
				,type: objectiveType
			},
			constraints: {
				rows: numConstraints,
				cols: numVariables + 1,
				vars: constraints.flat(),
				signs: constraintSigns
			}
		};

		try {
			const response = await fetch('/process', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(requestBody)
			});

			if (response.ok) {
				const data = await response.json();
				setResult(data);
			} else {
				setResult({ error: 'Error al resolver el problema' });
			}
		} catch (error) {
			setResult({ error: 'Error de conexión con el servidor' });
		}

		setLoading(false);
	};

	const downloadPDF = async () => {
		const requestBody = {
			objective: {
				n: numVariables,
				coefficients: objective,
				type: objectiveType
			},
			constraints: {
				rows: numConstraints,
				cols: numVariables + 1,
				vars: constraints.flat(),
				signs: constraintSigns
			}
		};

		try {
			const response = await fetch('/process?format=pdf', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(requestBody)
			});

			if (response.ok) {
				const blob = await response.blob();
				const url = window.URL.createObjectURL(blob);
				const a = document.createElement('a');
				a.href = url;
				a.download = 'resultado_simplex.pdf';
				document.body.appendChild(a);
				a.click();
				window.URL.revokeObjectURL(url);
				document.body.removeChild(a);
			} else {
				alert('Error al generar el PDF');
			}
		} catch (error) {
			alert('Error de conexión con el servidor');
		}
	};

	return (
		<div class="container">
			<h1>Autosimplex - Método Simplex</h1>
			
			<div class="config-section">
				<h2>Configuración del problema</h2>
				<div class="config-inputs">
					<label>
						Variables: 
						<input 
							type="number" 
							value={numVariables} 
							min="1" 
							max="100"
							step="1"
							onChange={(e) => handleVariablesChange((e.target as HTMLInputElement).value)}
						/>
					</label>
					<label>
						Restricciones: 
						<input 
							type="number" 
							value={numConstraints} 
							min="1" 
							max="100"
							step="1"
							onChange={(e) => handleConstraintsChange((e.target as HTMLInputElement).value)}
						/>
					</label>
				</div>
			</div>

			<div class="objective-section">
				<h2>Función objetivo ({objectiveType === 'maximize' ? 'maximizar' : 'minimizar'})</h2>
				<div class="objective-type">
					<label>
						<input type="radio" name="objectiveType" value="maximize" checked={objectiveType === 'maximize'} onChange={(e) => setObjectiveType((e.target as HTMLInputElement).value as 'maximize' | 'minimize')} /> Maximizar
					</label>
					<label>
						<input type="radio" name="objectiveType" value="minimize" checked={objectiveType === 'minimize'} onChange={(e) => setObjectiveType((e.target as HTMLInputElement).value as 'maximize' | 'minimize')} /> Minimizar
					</label>
				</div>
				<div class="objective-inputs">
					{objective.map((coeff, index) => (
						<label key={index}>
							x<sub>{index + 1}</sub>:
							<input
								type="number"
								step="0.01"
								value={coeff}
								onChange={(e) => updateObjective(index, (e.target as HTMLInputElement).value)}
							/>
						</label>
					))}
				</div>
			</div>

			<div class="constraints-section">
				<h2>Restricciones</h2>
				<div class="constraints-table">
					{constraints.map((constraint, rowIndex) => (
						<div key={rowIndex} class="constraint-row">
							<span class="constraint-label">Restricción {rowIndex + 1}:</span>
							{constraint.map((value, colIndex) => (
								<span key={colIndex} class="constraint-input">
									{colIndex < numVariables ? (
										<>
											<input
												type="number"
												step="0.01"
												value={value}
												onChange={(e) => updateConstraint(rowIndex, colIndex, (e.target as HTMLInputElement).value)}
											/>
											<span>x<sub>{colIndex + 1}</sub></span>
											{colIndex < numVariables - 1 && <span> + </span>}
										</>
									) : (
										<>
											<select value={constraintSigns[rowIndex]} onChange={(e) => {
												const newSigns = [...constraintSigns];
												newSigns[rowIndex] = (e.target as HTMLSelectElement).value;
												setConstraintSigns(newSigns);
											}}>
												<option value="<=">≤</option>
												<option value=">=">≥</option>
												<option value="=">=</option>
											</select>
											<input
												type="number"
												step="0.01"
												value={value}
												onChange={(e) => updateConstraint(rowIndex, colIndex, (e.target as HTMLInputElement).value)}
											/>
										</>
									)}
								</span>
							))}
						</div>
					))}
				</div>
			</div>

			<div class="solve-section">
				<button 
					class="solve-button" 
					onClick={solveSimplex}
					disabled={loading}
				>
					{loading ? 'Resolviendo...' : 'Resolver problema'}
				</button>
			</div>

			{result && (
				<div class="result-section">
					<h2>Resultado</h2>
					{result.error ? (
						<div class="error">{result.error}</div>
					) : (
<>
    {/* 1. SECCIÓN DE ADVERTENCIA (Tomada de la versión 'main') */}
    {result.warning && result.warning.trim() !== '' && (
        <div class="warning-box">
            <div class="warning-icon">⚠️</div>
            <div class="warning-content">
                <h3>Advertencia</h3>
                <p>{result.warning}</p>
            </div>
        </div>
    )}

    <div class="success">
        {/* 2. VALOR ÓPTIMO CONDICIONAL (Ajustado de la versión 'main') */}
        {!result.warning?.includes('infactible') && (
            <p><strong>Valor óptimo:</strong> {result.optimal_value}</p>
        )}

        {/* 3. ETIQUETA DE SOLUCIÓN CONDICIONAL (Ajustado de la versión 'main') */}
        <p><strong>{result.warning?.includes('infactible') ? 'Solución parcial alcanzada:' : 'Solución:'}</strong></p>

        {/* 4. LISTA DE SOLUCIÓN (Presente en ambas) */}
        <ul>
            {result.solution?.map((value, index) => (
                <li key={index}>x<sub>{index + 1}</sub> = {value.toFixed(4)}</li>
            ))}
        </ul>
        
        {/* 5. BOTÓN DE DESCARGA CONDICIONAL (Ajustado de la versión 'main') */}
        {!result.warning?.includes('infactible') && (
            <button
                class="download-pdf-button"
                onClick={downloadPDF}
            >
                📄 Descargar resultado en PDF
            </button>
        )}

        {/* 6. SECCIÓN DE TABLAS INTERMEDIAS (Tomada de la versión '57-Visualizacion-de-tablas-intermedias-Frontend') */}
        {result.steps && result.steps.length > 0 && (
            <div class="steps-section">
                <h3>Tablas intermedias</h3>
                {result.steps.map((step, idx) => {
                    const totalVars = step.cj ? step.cj.length : 0;
                    const n = result.solution ? result.solution.length : 0;
                    
                    // Construir nombres de variables: X1..Xn luego S1..Sextra
                    const varNames = [] as string[];
                    for (let i = 0; i < totalVars; i++) {
                        if (i < n) varNames.push(`X${i + 1}`);
                        else varNames.push(`S${i - n + 1}`);
                    }
                    
                    return (
                        <div key={idx} class="step-block">
                            <h4>Iteración {step.iteration}</h4>
                            {/* Fila Cj */}
                            <table class="simplex-table">
                                <thead>
                                    <tr class="cj-row">
                                        <th></th>
                                        <th></th>
                                        {(step.cj || []).map((v: number, j: number) => (
                                            <th key={j} class="cj-cell">{v.toFixed(2)}</th>
                                        ))}
                                        <th>R</th>
                                    </tr>
                                    <tr class="varnames-row">
                                        <th>c_b</th>
                                        <th>Base</th>
                                        {varNames.map((vn, j) => (
                                            <th key={j}>{vn}</th>
                                        ))}
                                        <th>R</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {(step.table || []).map((row: number[], rIdx: number) => (
                                        <tr key={rIdx}>
                                            <td class="cb-cell">{(step.cb && step.cb[rIdx] != null) ? step.cb[rIdx].toFixed(2) : ''}</td>
                                            <td class="base-cell">
                                                {(() => {
                                                const bv = step.base_variables?.[rIdx];
                                                if (!bv) return '';
                                                if (bv <= n) return `X${bv}`;
                                                return `S${bv - n}`;
                                            })()}
                                            </td>
                                            {row.slice(0, totalVars).map((cell: number, cIdx: number) => (
                                                <td key={cIdx} className={(step.pivot_row === rIdx && step.pivot_col === cIdx) ? 'pivot' : ''}>{cell.toFixed(2)}</td>
                                            ))}
                                            <td class="r-cell">{row[totalVars] != null ? row[totalVars].toFixed(2) : ''}</td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                            {/* Mensaje del Pivote */}
                            <div class="pivot-message">
                                Ingresa la variable <strong>{step.entering_var <= n ? `X${step.entering_var}` : `S${step.entering_var - n}`}</strong> y sale de la base la variable <strong>{step.leaving_var <= n ? `X${step.leaving_var}` : `S${step.leaving_var - n}`}</strong>. El elemento pivote es <strong>{(() => {
                                    const pv = (step.table && step.table[step.pivot_row]) ? step.table[step.pivot_row][step.pivot_col] : null;
                                    return pv != null ? pv.toFixed(2) : step.t_value?.toFixed?.(2) || '';
                                })()}</strong>
                            </div>
                        </div>
                    );
                })}
            </div>
        )}
    </div>
</>
					)}
				</div>
			)}
		</div>
	);
}

render(<App />, document.getElementById('app'));
