import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { fetchRandomChampion } from '../utils/api';

interface Champion {
    name: string;
    imageURL: string;
}

interface GameContextType {
    score: number;
    setScore: React.Dispatch<React.SetStateAction<number>>;
    champion: any;
    setChampion: React.Dispatch<React.SetStateAction<any>>;
    timeLeft: number;
    setTimeLeft: React.Dispatch<React.SetStateAction<number>>;
    timerRunning: boolean;
    setTimerRunning: React.Dispatch<React.SetStateAction<boolean>>;
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
        if (!timerRunning || timeLeft <= 0) {
            if (timeLeft <= 0) {
                setTimerRunning(false);
                setTimeLeft(allottedTime);
            }
            return;
        }

        const timerId = setInterval(() => {
            setTimeLeft((prev) => prev - 1);
        }, 1000);

        return () => clearInterval(timerId);
    }, [timerRunning, timeLeft]);
    
    useEffect(() => {
        const initializeChampion = async () => {
            const {name, imageURL} = await fetchRandomChampion();
            setChampion({name: name, imageURL: imageURL});
        };
    
        initializeChampion();
    }, []);
    

    const value = {
        score,
        setScore,
        champion,
        setChampion,
        timeLeft,
        setTimeLeft,
        timerRunning,
        setTimerRunning,
        allottedTime,
    };

    return <GameContext.Provider value={value}>{children}</GameContext.Provider>;
};

export const useGame = () => {
    const context = useContext(GameContext);
    if (!context) {
        throw new Error('useGame must be used within a GameProvider');
    }
    return context;
};
