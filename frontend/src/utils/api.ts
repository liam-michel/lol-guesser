import { json } from "stream/consumers";

// src/utils/api.ts
export async function fetchRandomChampion() {
  try {
      const port = import.meta.env.VITE_GOLANG_PORT;
      const fullURL = `http://localhost:${port}/api/randomchampion`;

      const response = await fetch(fullURL);
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

export async function createUser(username: string, password: string){
    try{
        const port = import.meta.env.VITE_GOLANG_PORT;
        const fullURL = `http://localhost:${port}/api/createuser`;

        const response = await fetch(fullURL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({username: username, password: password}),
        });

        if (!response.ok) {
            throw new Error("HTTP error when creating new user! Status: " + response.status);
        }

        const data = await response.json();
        return data;
    }catch(error){
        console.error("Error creating new user: ", error);
        return JSON.stringify({error: error});
    }
}