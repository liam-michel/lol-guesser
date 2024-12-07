import {Champion} from "../types";
import { fetchRandomChampion } from '../utils/api';
export default function RandomChampion({champion, setChampion}: {champion: Champion | null, setChampion: any}){
    //call GO backend to get a random champion
    const handleFetchRandomChampion = async () => {
        const {name, imageURL} = await fetchRandomChampion();
        setChampion({name, imageURL})
    }

    return(
        <div>
            <h1>Fetch new random champion</h1>
            <button onClick={handleFetchRandomChampion}>Fetch new Champion</button>
            <h2>{champion?.name}</h2>
            <img src={champion?.imageURL || ''} alt={champion?.name || ''} />
        </div>
    )
}