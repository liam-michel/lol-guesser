import { useGame } from "../context/GameContext";
export default function Score() {
    const { accuracy, correctGuesses, incorrectGuesses } = useGame();


    return(
    <>
    <div className = "metricscontainer">
        <h2>Correct Guesses: {correctGuesses}</h2>
        <h2>Incorrect Guesses: {incorrectGuesses}</h2>
        <h2>Accuracy: {accuracy}%</h2>
    </div>
    </>
    )
}
