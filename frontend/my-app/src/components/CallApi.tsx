import React , {useState, useEffect} from 'react';

export default function CallApi(){
    const [message, setMessage] = useState("");

    const fetchMessage = async() => {
        try{
            const response = await fetch("http://localhost:8080/testAPI")
            if(!response.ok){
                throw new Error("HTTP error! status: " + response.status)
            }
            const data = await response.json()
            setMessage(data.message)
        } catch (error){
            console.error("Error fetching message: ", error)
            setMessage("Error fetching message")
        }
    }

    return(
        <div>
            <button onClick={fetchMessage}>Fetch Message</button>
            <p>{message}</p>
        </div>
    )
}