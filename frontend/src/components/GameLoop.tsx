import { useEffect, useState } from "react";
import { useGame } from "../context/GameContext";  // Import the useGame hook
import TitleBar from "./TitleBar";
import RandomChampion from "./RandomChampion";

export default function GameLoop() {
    const stateValues = useGame();
    const {
        score = 0,
        setScore,
        champion,
        setChampion,
        timeLeft = 0,
        setTimeLeft,
        timerRunning = false,
        setTimerRunning,
        allottedTime = 0
    } = stateValues || {};

    const [inputValue, setInputValue] = useState(''); // Keeping local state for input value

    // Start timer when the game starts and when timeLeft changes
    useEffect(() => {
        if (!timerRunning || timeLeft <= 0) {
            if (timeLeft <= 0) {
                setTimerRunning(false); // Stop the timer
                setTimeLeft(allottedTime); // Reset the timer
            }
            return;
        }

        // Start the timer
        const timerId = setInterval(() => {
            setTimeLeft((prevTime) => prevTime - 1);
        }, 1000);

        // Cleanup the timer when the component unmounts or dependencies change
        return () => clearInterval(timerId);
    }, [timerRunning, timeLeft, allottedTime, setTimeLeft]);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setInputValue(e.target.value);
        if (!timerRunning) {
            setTimerRunning(true); // Start the timer if it's not running
        }
    };

    const checkGuess = () => {
        console.log("Checking guess: ", inputValue);
        const formattedGuess = inputValue.toLowerCase().replace(/\s+/g, '');
        const formattedChampionName = champion?.name.toLowerCase().replace(/\s+/g, '');
        if (formattedGuess === formattedChampionName) {
            console.log("Correct guess!");
            setScore(score + 1);  // Increment the score
            // Optionally, you can fetch a new champion here after a correct guess.
            // setChampion(someNewChampion); 
        } else {
            console.log("Incorrect guess!");
        }
    };

    return (
        <div className="game-loop">
            <TitleBar />
            <RandomChampion champion={champion} setChampion={setChampion}/>
            <div className="stats">
                <input
                    type="text"
                    value={inputValue}
                    onChange={handleInputChange}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                            checkGuess();
                        }
                    }}
                    style={{
                        padding: '10px',
                        fontSize: '20px',
                        width: '300px',
                        height: '40px',
                    }}
                />
                <h1>Score: {score}</h1>
                <h1>Timer: {timeLeft}</h1>
            </div>
        </div>
    );
}
