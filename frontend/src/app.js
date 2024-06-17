import React from 'react';
import AppRoutes from './routes';

// Компонент App является основным компонентом приложения
function App() {
    return (
        <div>
            {/* Используем маршруты, определенные в AppRoutes */}
            <AppRoutes />
        </div>
    );
}

export default App;
