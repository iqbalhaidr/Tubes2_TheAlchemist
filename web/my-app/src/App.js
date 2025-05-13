import { useState } from 'react';
import { Network } from 'vis-network/standalone';
import 'vis-network/styles/vis-network.css';
import './App.css';

// Program Aplikasi berbasis Web
function App() {
  const [algorithm, setAlgorithm] = useState('bfs');
  const [search, setSearch] = useState('satu');
  const [element, setElement] = useState(' '); 
  const [searchTime, setSearchTime] = useState('00.00');
  const [nodesVisited, setNodesVisited] = useState('0');
  const [recipeResult, setRecipeResult] = useState('Masukkan nama element yang ingin dicari, lalu klik tombol "Cari Resep"');
  const [maxRecipes, setMaxRecipes] = useState(2);

    // Fetch data API
    const findRecipes = async () => {
        if (!element.trim()) {
            alert('Tolong masukkan nama element yang benar!');
            return;
        }

        // const startTime = performance.now();

        try {
          const query = new URLSearchParams({
            element : element.trim(), 
            multiple: search === 'semua' ? true : false,
            n: search === 'semua' ? maxRecipes : 1
          }).toString();

          // Memilih algoritma
          const endpoint = algorithm === 'bfs' ? '/bfs' : '/dfs';
          const responseData = await fetch(`http://localhost:8080${endpoint}?${query}`,{
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            },
          });

          if (!responseData.ok) {
            const errorText = await responseData.text();
            throw new Error(`Gagal memuat data dari backend: ${responseData.status} ${responseData.statusText} - ${errorText}`);
          }

          const responseDataJson = await responseData.json();
          console.log('Backend response: ' , responseDataJson);

          // Penyesuaian Format Backend ke FrontEnd
          let searchTime = responseDataJson.SearchTimeInMiliseconds || 0;
          let nodesVisited = (responseDataJson.NodeCountElement || 0) + (responseDataJson.NodeCountRecipe || 0);
          

          let nodes = [];
          let edges = [];
          let idCounter = 1;
          
          
          if (responseDataJson.isFound !== undefined) {
            
            // Satu Resep
            const steps = responseDataJson.Steps || [];
            

            if (!steps || steps.length === 0) {
                setRecipeResult(`Tidak ada resep yang ditemukan!`);
                return;
            }
            
            const rootId = idCounter++;

            let prevResultId = rootId;
            for (let i = 0; i < steps.length; i++) {
                
                const s = steps[i];
                const id1 = idCounter++;
                const id2 = idCounter++;
                const idResult = idCounter++;

                // Tambah Simpul/Nodes
                nodes.push({ id: id1, label: s[0] });
                nodes.push({ id: id2, label: s[1] });
                nodes.push({ id: idResult, label: s[2] });

                // Tambah Sisi/Edges
                edges.push({ from: id1, to: idResult });
                edges.push({ from: id2, to: idResult });

                // Menghubungkan Simpul
                if ( i === 0) {
                    edges.push({ from: rootId , to: id1 });
                    edges.push({ from: rootId , to: id2 });
                } else {
                    const prevStep = steps[i-1];
                    if (prevStep[2] === s[0]) {
                        edges.push({ from: prevResultId, to: id1 }); 
                    } else if (prevStep[2] === s[1]) {
                        edges.push({ from: prevResultId, to: id2 });
                    } 
                    else {
                        edges.push({ from: rootId, to: id1 });
                        edges.push({ from: rootId, to: id2 });
                    }
                }
                prevResultId = idResult;
            }
            
                
           } else if (responseDataJson.isSatisfied !== undefined) {
            
            // Resep Banyak/Multiple
            const paths = responseDataJson.Steps || [];
            let idCounter = 1;

            if (!paths || paths.length === 0) {
                setRecipeResult(`Tidak ada resep yang ditemukan!`);
                return;
            }

            for (let i = 0; i <paths.length; i++) {
                const pathSteps = paths[i];
                let pathNodes = []
                let pathEdges = []

                for (let j = 0; j < pathSteps.length; j++) {
                    const s = pathSteps[j];

                    const id1 = idCounter++;
                    const id2 = idCounter++;
                    const idResult = idCounter++;

                    // Tambah Simpul/Nodes
                    pathNodes.push({ id: id1, label: s[0], group: i });
                    pathNodes.push({ id: id2, label: s[1], group: i });
                    pathNodes.push({ id: idResult, label: s[2], group: i });

                    // Tambah Sisi/Edges
                    pathEdges.push({ from: id1, to: idResult });
                    pathEdges.push({ from: id2, to: idResult });

                    // Menghubungkan Simpul
                    if (j < pathSteps.length - 1) {
                        const nextStep = pathSteps[j+1];

                        if (nextStep[0] === s[2]) {
                            const nextId1 = idCounter;
                            pathEdges.push({ from: idResult, to: nextId1 });
                        } else if (nextStep[1] === s[2]) {
                            const nextId2 = idCounter + 1;
                            pathEdges.push({ from: idResult, to: nextId2})
                        }
                    }
                    
                }

                nodes = nodes.concat(pathNodes);
                edges = edges.concat(pathEdges);
            }
           }
           
        

          setSearchTime(searchTime.toFixed(2) || ' N/A');
          setNodesVisited(nodesVisited || ' N/A');

          let recipeCount = 0;
          if (responseDataJson.isFound !== undefined) {
              recipeCount = 1;
          } else if (responseDataJson.isSatisfied !== undefined) {
              recipeCount = (responseDataJson.Steps || []).length || 0;
          }
          setRecipeResult(search === 'satu' ? `Ditemukan satu resep untuk ${element}` : `Ditemukan ${recipeCount} resep untuk ${element}`);


          const container = document.getElementById('tree');
          const data = { nodes, edges };
          const options = {
                layout: { 
                    hierarchical: { 
                        direction: 'UD', 
                        sortMethod: 'directed',
                        
                        nodeSpacing: 100,
                        levelSeparation: 50,
                        treeSpacing: 200
                    }}, 
                edges: { arrows: 'to' },
                nodes: {
                    shape: 'box',
                    font: { size: 12},
                    color: {
                        background: 'lightgray',
                        border: 'black'
                    },
                },
                physics: {
                    enabled: true
                },
                interaction: {
                    hover: true
                }
            };
            new Network(container, data, options);
        } catch(error) {
            console.error('Error:', error);
            setRecipeResult(`Gagal mengambil resep: ${error.message}`)
        }
      };
    


    // Tampilan User Interface atau UI
    return (
        <div style={{ padding: '20px' }}>
            <header className="App-header">
                <h1 className="App-title">Tugas Besar 2</h1>
                <h2 className="App-subtitle">Pencarian Resep dengan Algoritma BFS/DFS dalam Permainan Little Alchemy 2</h2>
                <h3 className="App-subsubtitle">Dibuat oleh Kelompok The Alchemist</h3>
            </header>

            <div>
                <label>Algoritma: </label>
                <select value={algorithm} onChange={(e) => setAlgorithm(e.target.value)}>
                    <option value="bfs">BFS</option>
                    <option value="dfs">DFS</option>
                </select>
            </div>

            <div>
                <label>Pencarian: </label>
                <select value={search} onChange={(e) => setSearch(e.target.value)}>
                    <option value="resep">Satu Resep</option>
                    <option value="semua">Semua Resep</option>
                </select>
                {search === 'semua' && (
                    <input
                        type="number"
                        value={maxRecipes}
                        onChange={(e) => setMaxRecipes(Math.max(1, parseInt(e.target.value)))}
                        min="1"
                        placeholder="Jumlah Resep Terbanyak"
                    />
                )}
            </div>

            <div>
                <label>Nama Element: </label>
                <input 
                type="text" 
                value={element} 
                onChange={(e) => setElement(e.target.value)}
                placeholder="Element"
                />
                <button onClick={findRecipes}>Cari Resep</button>
            </div>

            <div className="App-output">
                <div>
                    <p>Waktu Pencarian: {searchTime} ms</p>
                    <p>Jumlah Node yang Dikunjungi: {nodesVisited}</p>
                    <p>Hasil Pencarian: {recipeResult}</p>
                </div> 
                <div
                id="tree"
                style={{
                    height: '500px',
                    width: '100%',
                    border: '1px solid lightgray',
                }}
            />
        </div>
    </div>
    )
}

export default App;
