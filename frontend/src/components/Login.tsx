import { useState } from 'react';
import { useAuth } from '../context/AuthContext';

export default function Login() {
  const [inputName, setInputName] = useState('');
  const [inputPassword, setInputPassword] = useState('');
  const [error, setError] = useState('');
  const { username, token, login, createAccount, logout, isAuthenticated } = useAuth();

  // Generic function to handle both login and create account
  const handleSubmit = async (event: any, action: 'login' | 'createAccount') => {
    event.preventDefault();
    setError(''); // Reset error before making API call
  
    if (inputName === '' || inputPassword === '') {
      setError('Please fill in both fields.');
      return;
    }
  
    try {
      let response;
      if (action === 'createAccount') {
        response = await createAccount(inputName, inputPassword);
      } else if (action === 'login') {
        response = await login(inputName, inputPassword);
      }
  
      if (response?.error) {
        setError(response.error); // Handle the error and display it
      } else {
        // Success: You might want to do something else (e.g., redirect, show success message)
        setError(''); // Reset error if successful
      }
    } catch (error) {
      setError("An unexpected error occurred"); // Generic error handling
    }
  };

  return (
    <div className="login-container">
      {isAuthenticated ? (
        // Render username and token if authenticated
        <div>
          <h2>Welcome, {username ? username : 'Guest'}!</h2>
          <button onClick={logout}>Logout</button>
        </div>
      ) : (
        // Render login or sign up form if not authenticated
        <form>
          <input
            type="text"
            placeholder="Username"
            onChange={(e) => setInputName(e.target.value)}
            value={inputName}
          />
          <input
            type="password"
            placeholder="Password"
            onChange={(e) => setInputPassword(e.target.value)}
            value={inputPassword}
          />
          <button onClick={(e) => handleSubmit(e, 'createAccount')}>Create Account</button>
          <button onClick={(e) => handleSubmit(e, 'login')}>Login</button>
          {error && <p style={{ color: 'red' }}>{error}</p>}
        </form>
      )}
    </div>
  );
}
