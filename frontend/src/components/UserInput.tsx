import React, { useState } from "react";
import { useGame } from "../context/GameContext";
import { fetchRandomChampion } from "../utils/api";
export default function UserInput() {
    const { champion, setChampion, timerRunning, startTimer, inputValue, setInputValue, correctGuesses, setCorrectGuesses, incorrectGuesses, setIncorrectGuesses, accuracy, setAccuracy } = useGame();

    const updateAccuracy = () => {
        //calculate accuracy as the number of correct guesses divided by the total number of guesses
        const totalGuesses = correctGuesses + incorrectGuesses
        const accuracy = (correctGuesses / totalGuesses) * 100;
        const rounded = parseFloat(accuracy.toFixed(2));
        if (isNaN(rounded)) {
            setAccuracy(0);
            return;
        }
        console.log("new accuracy", rounded);   
        setAccuracy(rounded);
    }

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

            setCorrectGuesses(correctGuesses + 1);
        }
        else{
            setIncorrectGuesses(incorrectGuesses + 1);
        }
        //update statistics
    
        updateAccuracy();
        //empty the user input ready for another one      
        setInputValue(''); // Reset input after each guess
    };

    const handleKeyDown = async (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            checkGuess();
            //fetch a new champion
            const { name, imageURL } = await fetchRandomChampion();
            setChampion({ name, imageURL });
            setInputValue('');
        }
    };

    return (
        <>
        <h2>Enter Champion Name</h2>>
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
                marginBottom: '20px'
            }}
        />

        </>
    );
}
