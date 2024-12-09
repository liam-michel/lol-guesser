import { useGame } from "../context/GameContext";
// import { fetchRandomChampion } from "../utils/api";

export default function RandomChampionDisplay() {
    const { champion, setChampion } = useGame(); // Use context to manage state

    // const handleFetchChampion = async () => {
    //     try {
    //         const newChampion = await fetchRandomChampion();
    //         setChampion(newChampion); // Update the context with the new champion
    //     } catch (error) {
    //         console.error("Failed to fetch new champion:", error);
    //     }
    // };

    return (
        <div>
            <img src={champion?.imageURL || ''} alt={champion?.name || ''} style = {{
                width: "250px",
                height: "250px",
                borderRadius: "50%",
                border: "2px solid black",
                boxShadow: "5px 5px 5px #000000"
            }} />
        </div>
    );
}
