import React, {useState, useEffect} from 'react';

export default function RandomChampion(){
    const  [champion, setChampion] = useState("");
    const [image, setImage] = useState(null);
    
    //call GO backend to get a random champion

    const fetchRandomChampion = async() => {
        try{
            const response = await fetch("http://localhost:8080/api/randomchampion")
            if(!response.ok){
                throw new Error("HTTP error when fetching new random champopipn! status: " + response.status)
            
            }
            const data = await response.json()

        }
    }
    return(
        <div>
            <h1>Fetch new random champion</h1>
        </div>
    )
}