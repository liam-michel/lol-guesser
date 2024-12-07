import { useGame } from "../context/GameContext";
import { fetchRandomChampion } from "../utils/api";

export default function RandomChampionDisplay() {
    const { champion, setChampion } = useGame(); // Use context to manage state

    const handleFetchChampion = async () => {
        try {
            const newChampion = await fetchRandomChampion();
            setChampion(newChampion); // Update the context with the new champion
        } catch (error) {
            console.error("Failed to fetch new champion:", error);
        }
    };

    return (
        <div>
            <h1>Fetch new random champion</h1>
            <button onClick={handleFetchChampion}>Fetch new Champion</button>
            <h2>{champion?.name}</h2>
            <img src={champion?.imageURL || ''} alt={champion?.name || ''} />
        </div>
    );
}
