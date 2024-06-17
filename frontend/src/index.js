import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import './index.css';

// Рендерим компонент App в элемент с id 'root' в файле index.html
ReactDOM.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>,
    document.getElementById('root')
);
