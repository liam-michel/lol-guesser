import React from "react";
import { useGame } from "../context/GameContext";

export default function Score() {
    const { score } = useGame();

    return <h1>Score: {score}</h1>;
}
