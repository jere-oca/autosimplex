import { render } from 'preact';
import { useState } from 'preact/hooks';
import './style.css';

export function App() {
	const [numVariables, setNumVariables] = useState(2);
	const [numConstraints, setNumConstraints] = useState(2);
	const [objective, setObjective] = useState([1, 1]);
	const [objectiveType, setObjectiveType] = useState('maximize');
	const [constraints, setConstraints] = useState([[1, 1, 1], [1, 2, 2]]);
	const [constraintSigns, setConstraintSigns] = useState<string[]>(Array(2).fill("<="));
	const [result, setResult] = useState(null);
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
						<input type="radio" name="objectiveType" value="maximize" checked={objectiveType === 'maximize'} onChange={(e) => setObjectiveType((e.target as HTMLInputElement).value)} /> Maximizar
					</label>
					<label>
						<input type="radio" name="objectiveType" value="minimize" checked={objectiveType === 'minimize'} onChange={(e) => setObjectiveType((e.target as HTMLInputElement).value)} /> Minimizar
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
						<div class="success">
							<p><strong>Valor óptimo:</strong> {result.optimal_value}</p>
							<p><strong>Solución:</strong></p>
							<ul>
								{result.solution?.map((value, index) => (
									<li key={index}>x<sub>{index + 1}</sub> = {value.toFixed(4)}</li>
								))}
							</ul>
						</div>
					)}
				</div>
			)}
		</div>
	);
}

render(<App />, document.getElementById('app'));
