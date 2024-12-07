import React, {useState, useEffect} from 'react';

export default function RandomChampion(){
    const  [champion, setChampion] = useState("");
    const [image, setImage] = useState("");
    
    //call GO backend to get a random champion

    const fetchRandomChampion = async() => {
        try{
            const port = import.meta.env.VITE_GOLANG_PORT
            console.log("port ", port)
            const fullURL = `http://localhost:${import.meta.env.VITE_GOLANG_PORT}/api/randomchampion`
            console.log("Full URL: ", fullURL)
            const response = await fetch(fullURL)
            if(!response.ok){
                throw new Error("HTTP error when fetching new random champopipn! status: " + response.status)
            
            }
            const data = await response.json()
            setChampion(data.name)
            setImage(data.url)
            
        }catch(error){
            console.error("Error fetching random champion: ", error)
        }
    }
    return(
        <div>
            <h1>Fetch new random champion</h1>
            <button onClick={fetchRandomChampion}>Fetch new Champion</button>
            <h2>{champion}</h2>
            <img src={image} alt={champion} />
        </div>
    )
}