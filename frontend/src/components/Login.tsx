import {useState, useEffect} from 'react';
import { createUser } from '../utils/api';

export default function Login() {
  const [username, setUsername] = useState(''); 
  const [password, setPassword] = useState('');
  const [response, setResponse] = useState('');
  const [error, setError] = useState('');

  const attemptCreateUser = async (event: any) => {
    event?.preventDefault()
    console.log("Username: ", username);
    console.log("Password: ", password);
    if(username === '' || password === '') {
      return
    };
    try {
      const response = await createUser(username, password);
      console.log("Received response from attempting to create user: ", response);
      if (response.error) {
        setError(response.error);
      } else {
        setResponse(response);
        // Redirect to the home page
      }
    } catch (error) {
      setError
    }
  }
  
  return (
    <div className="login-container">
        <form>
            <input type="text" placeholder="Username" onChange = {(e) => setUsername(e.target.value)} />
            <input type="password" placeholder="Password" onChange = {(e) => setPassword(e.target.value)} />
            <button onClick = {attemptCreateUser}>Create Account</button>
            <button type="submit">Login</button>
        </form>
    </div>
);
}