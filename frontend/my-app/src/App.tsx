import React from 'react';
import logo from './logo.svg';
import './App.css';
import CallApi from './components/CallApi';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <CallApi />
      </header>
    </div>
  );
}

export default App;
