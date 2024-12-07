import {useState} from 'react';

export default function CallApi(){
    const [message, setMessage] = useState("");

    const fetchMessage = async() => {
        try{
            const fullURL = `http://localhost:${import.meta.env.VITE_GOLANG_PORT}/api/testAPI`
            const response = await fetch(fullURL)
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