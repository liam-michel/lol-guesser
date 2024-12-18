import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import App from './App.tsx';
import { GameProvider }  from './context/GameContext.tsx';
import { AuthProvider } from './context/AuthContext.tsx';
createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider>
      <GameProvider> {/* Wrap the App with the GameProvider */}
        <App />
      </GameProvider>
    </AuthProvider>
    </StrictMode>,
);
