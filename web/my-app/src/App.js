const React = require('react');
const { useState } = React;
const { Network } = require('vis-network/standalone');
require('vis-network/styles/vis-network.css');

// Program Aplikasi berbasis Web
function App() {
  const [algorithm, setAlgorithm] = useState('bfs');
  const [search, setSearch] = useState('terpendek');
  const [material, setMaterial] = useState(' '); 
  const [searchTime, setSearchTime] = useState('00.00');
  const [nodesVisited, setNodesVisited] = useState('0');
  const [recipeResult, setRecipeResult] = useState('Masukkan nama material yang ingin dicari, lalu klik tombol "Cari Resep"');
  const [maxRecipes, setMaxRecipes] = useState(2);

    // Fetch data API
    const findRecipes = async () => {
        if (!material.trim()) {
            alert('Tolong masukkan nama material yang benar!');
            return;
        }
        // const startTime = performance.now();

        try {
          const responseData = await fetch('http://localhost:8000/recipes', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              material: material.trim(),
              algorithm: algorithm,
              search: search,
              maxRecipes: maxRecipes,
            }),
          });

          if (!responseData.ok) {
            throw new Error('Gagal memuat data dari backend');
          }

          const dataResponse = await responseData.json();
          console.log('Backend response: ' , dataResponse);

          setSearchTime(dataResponse.searchTime.toFixed(2) || ' N/A');
          setNodesVisited(dataResponse.nodesVisited || ' N/A');
          setRecipeResult(search === 'terpendek' ? `Resep terpendek untuk ${material}` : `Ditemukan ${recipes.length / 2} resep untuk ${material}`);

          const recipes = dataResponse.recipes;
          const nodes = recipes.filter(item => item.label).map(item => ({ id: item.id, label: item.label }));
          const edges = recipes.filter(item => item.from).map(item => ({ from: item.from, to: item.to }));
        

        // const startTime = performance.now();
        // const visitedNodes = Math.floor(Math.random() * 100);
        // let recipes = [];

        // if (search === 'terpendek') {
        //     recipes = [
        //         { id: 1, label: material },
        //         { id: 2, label: 'fire' },
        //         { id: 3, label: 'earth' },
        //         { from: 2, to: 1 },
        //         { from: 3, to: 1 },
        //     ]
        // } else {
        //     recipes = [
        //         { id: 1, label: material },
        //         { id: 2, label: 'fire' },
        //         { id: 3, label: 'earth' },
        //         { id: 4, label: 'water' },
        //         { id: 5, label: 'air' },
        //         { from: 2, to: 1 },
        //         { from: 3, to: 1 },
        //         { from: 4, to: 1 },
        //         { from: 5, to: 1 },
        //     ].slice(0, maxRecipes * 2);
        // }
        
        // const endTime = performance.now();
        // setSearchTime((endTime - startTime).toFixed(2));
        // setNodesVisited(visitedNodes);
        // setRecipeResult(search === 'terpendek' ? `Resep terpendek untuk ${material}` : `Ditemukan ${recipes.length / 2} resep untuk ${material}`);


        // Visualisasi Tree
        // const nodes = recipes.filter(item => item.label).map(item => ({ id: item.id, label: item.label }));
        // const edges = recipes.filter(item => item.from).map(item => ({ from: item.from, to: item.to }));

        const container = document.getElementById('tree');
        const data = { 
          nodes, edges
        };
        const options = {
            layout: {hierarchical: { direction: 'UD', sortMethod: 'directed' }}, 
            edges: {arrows: 'to' }
        };
        const network = new Network(container, data, options);
      } catch (error) {
        console.error('Error:', error);
      }

    }


    // Tampilan User Interface atau UI
    return React.createElement(
        'div',
        { style: {padding: '20px' } },
        React.createElement(
            'header',
            { className: 'App-header' },
            React.createElement(
                'h1',
                { className: 'App-title' },
                'Tugas Besar 2'
            ),
            React.createElement(
              'h2',
              { className: 'App-subtitle' },
              'Pencarian Resep dengan Algoritma BFS/DFS dalam Permainan Little Alchemy 2'
            ),
            React.createElement(
              'h3',
              { className: 'App-subsubtitle'},
              'Dibuat oleh Kelompok The Alchemist'
            )
        ),
        React.createElement(
            'div',
            null,
            React.createElement(
                'label', null, 'Algoritma: '
            )
        ),
        React.createElement(
            'select',
            { value: algorithm, onChange: e => setAlgorithm(e.target.value) },
            React.createElement(
                'option', 
                { value: 'bfs' }, 'BFS'
            ),
            React.createElement(
                'option', 
                { value: 'dfs' }, 'DFS'
            ),
            React.createElement(
                'option', 
                { value: 'bidirectional' }, 'Bidirectional'
            )
        ),
        React.createElement(
            'div',
            null,
            React.createElement(
                'label', null, 'Pencarian: '
            ),
            React.createElement(
            'select',
            { value: search, onChange: e => setSearch(e.target.value) },
            React.createElement(
                'option', 
                { value: 'terpendek' }, 'Resep Terpendek'
            ),
            React.createElement(
                'option',
                {value: 'semua'}, 'Semua Resep')
            ),
            search === 'semua' &&
                React.createElement('input', {
                    type: 'number',
                    value: maxRecipes,
                    onChange: e => setMaxRecipes(Math.max(1, e.target.value)),
                    min: '1',
                    placeholder: 'Jumlah Resep Terbanyak'
                }) 
        ),
        React.createElement(
            'div',
            null,
            React.createElement(
                'label',
                null, 
                'Nama Material: '
            ),
            React.createElement(
                'input', 
                {
                    type: 'text',
                    value: material, 
                    onChange: e => setMaterial(e.target.value),
                    placeholder: 'e.g., water'
                }
            ),
            React.createElement(
                'button',
                { onClick: findRecipes },
                'Cari Resep'
            ),
        ),

        React.createElement(
            'div',
            { className: 'App-output' },
            React.createElement(
                'div',
                null,
                React.createElement(
                    'p',
                    null,
                    'Waktu pencarian: ',
                    searchTime,
                    ' ms'
                ),
                React.createElement(
                    'p',
                    null,
                    'Jumlah node yang dikunjungi: ',
                    nodesVisited
                ),
                React.createElement(
                    'p',
                    null,
                    'Hasil pencarian: ', 
                    recipeResult
                )
            ),
            React.createElement(
                'div',
                { id: 'tree',
                    style : 
                    { height: '400px', width: '100%', border: '1px solid #ccc', marginTop: '20px'   } 
                }
            )
        )
    )
}


module.exports = App;