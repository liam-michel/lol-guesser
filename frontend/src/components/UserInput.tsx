import React, { useState } from "react";
import { useGame } from "../context/GameContext";

export default function UserInput() {
    const { setScore, champion, score, timerRunning, startTimer } = useGame();
    const [inputValue, setInputValue] = useState('');

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setInputValue(e.target.value);
        if (!timerRunning) {
            startTimer(); // Start the timer using the GameContext method
        }
    };

    const checkGuess = () => {
        const formattedGuess = inputValue.toLowerCase().replace(/\s+/g, '');
        const formattedChampionName = champion?.name.toLowerCase().replace(/\s+/g, '');
        if (formattedGuess === formattedChampionName) {
            setScore(score + 1);
        }
        setInputValue(''); // Reset input after each guess
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            checkGuess();
        }
    };

    return (
        <input
            type="text"
            value={inputValue}
            onChange={handleInputChange}
            onKeyDown={handleKeyDown}
            style={{
                padding: '10px',
                fontSize: '20px',
                width: '300px',
                height: '40px',
            }}
        />
    );
}
