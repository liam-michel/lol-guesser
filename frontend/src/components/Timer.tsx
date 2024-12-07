import React from "react";
import { useGame } from "../context/GameContext";

export default function Timer() {
    const { timeLeft, startTimer, stopTimer, resetTimer } = useGame();

    return (
        <div>
            <h1>Timer: {timeLeft}</h1>
            <button onClick={startTimer}>Start</button>
            <button onClick={stopTimer}>Stop</button>
            <button onClick={resetTimer}>Reset</button>
        </div>
    );
}
