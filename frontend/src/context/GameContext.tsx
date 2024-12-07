import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { fetchRandomChampion } from "../utils/api";

interface Champion {
    name: string;
    imageURL: string;
}

interface GameContextType {
    score: number;
    setScore: React.Dispatch<React.SetStateAction<number>>;
    champion: Champion | null;
    setChampion: React.Dispatch<React.SetStateAction<Champion | null>>;
    timeLeft: number;
    startTimer: () => void;
    stopTimer: () => void;
    resetTimer: () => void;
    timerRunning: boolean;
    allottedTime: number;
}

const GameContext = createContext<GameContextType>({} as GameContextType);

export const GameProvider = ({ children }: { children: ReactNode }) => {
    const allottedTime = 30;
    const [score, setScore] = useState(0);
    const [champion, setChampion] = useState<Champion | null>(null);
    const [timeLeft, setTimeLeft] = useState(allottedTime);
    const [timerRunning, setTimerRunning] = useState(false);

    useEffect(() => {
        let timerId: NodeJS.Timeout | null = null;

        if (timerRunning && timeLeft > 0) {
            timerId = setInterval(() => {
                setTimeLeft((prevTime) => prevTime - 1);
            }, 1000);
        } else if (timeLeft <= 0) {
            setTimerRunning(false);
            setTimeLeft(allottedTime); // Reset the timer
        }

        return () => {
            if (timerId) clearInterval(timerId);
        };
    }, [timerRunning, timeLeft, allottedTime]);

    useEffect(() => {
        const initializeChampion = async () => {
            const { name, imageURL } = await fetchRandomChampion();
            setChampion({ name, imageURL });
        };

        initializeChampion();
    }, []);

    const startTimer = () => setTimerRunning(true);
    const stopTimer = () => setTimerRunning(false);
    const resetTimer = () => {
        setTimerRunning(false);
        setTimeLeft(allottedTime);
    };

    const value = {
        score,
        setScore,
        champion,
        setChampion,
        timeLeft,
        startTimer,
        stopTimer,
        resetTimer,
        timerRunning,
        allottedTime,
    };

    return <GameContext.Provider value={value}>{children}</GameContext.Provider>;
};

export const useGame = () => {
    const context = useContext(GameContext);
    if (!context) {
        throw new Error("useGame must be used within a GameProvider");
    }
    return context;
};
