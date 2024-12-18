import { customFetch } from './customFetch';

// src/utils/api.ts
const port = import.meta.env.VITE_GOLANG_PORT;

export async function fetchRandomChampion() {
  try {
      const fullURL = `http://localhost:${port}/api/randomchampion`;

      const response = await customFetch(fullURL);
      if (!response.ok) {
          throw new Error("HTTP error when fetching new random champion! Status: " + response.status);
      }

      const data = await response.json();
      return { name: data.name, imageURL: data.url }; // Return the champion object
  } catch (error) {
      console.error("Error fetching random champion: ", error);
      throw error; // Re-throw the error for caller to handle if needed
  }
}

export async function createUser(username: string, password: string) {
    try {
      const fullURL = `http://localhost:${port}/api/createuser`;
  
      const response = await customFetch(fullURL, {
        credentials: 'include',
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });
  
      if (!response.ok) {
        const errordata = await response.json();
        throw new Error(errordata.error); // Extract the error message and throw an Error
      }
  
      const data = await response.json();
      return data;
    } catch (error) {
      console.error("Error creating new user: ", error);
      // Return a clean error message, not the whole error object
      return { error: error instanceof Error ? error.message : "Unknown error occurred" };
    }
  }

export async function loginUser(username: string, password: string){
    try{
        const fullURL = `http://localhost:${port}/api/login`;

        const response = await customFetch(fullURL, {
            credentials: 'include',
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({username: username, password: password}),
        });

        if (!response.ok) {
            const errordata = await response.json();
            throw new Error(errordata.error);
        }

        const data = await response.json();
        return data;

    }catch(error){
        console.error("Error logging in: ", error);
        return { error: error instanceof Error ? error.message : "Unknown error occurred" };
    }
}

export function decodeJWT(token: string){
    try{
        const parts = token.split('.');
        if (parts.length !==3){
            return JSON.stringify({error: "Invalid JWT"});
        }
        const payload = parts[1];
        const decodedPayload = JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')));
        return decodedPayload;

    }catch(error){
        console.error("Failed to decode JWT: ", error);
        return JSON.stringify({error: error});  
    }
}

export function getTokenExpiry(token: string){
    try{
        const decoded = decodeJWT(token);
        if (decoded.exp){
            return new Date(decoded.exp * 1000);
        }else{
            return JSON.stringify({error: "No expiry date found in token"});
        }

    }catch(error){
        console.error("Failed to get token expiry: ", error);
        return JSON.stringify({error: error});
    }
}

export function isTokenExpired(token: string){
    try{
        const expiry = getTokenExpiry(token);
        return new Date() > expiry;

    }catch(error){
        console.error("Failed to check if token is expired: ", error);
        return JSON.stringify({error: error});
    }
}