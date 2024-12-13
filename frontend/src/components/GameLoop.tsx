import TitleBar from "./TitleBar";
import Timer from "./Timer";
import Metrics from "./Metrics";
import UserInput from "./UserInput";
import RandomChampion from "./RandomChampion";
import Login from "./Login";

export default function GameLoop() {
    return (
        <div className="game-loop">
            <div className="login-container">
                <Login />
            </div>
            <TitleBar />
            <div className="main-container">
                <div className="game-container">
                    <RandomChampion />
                    <UserInput />
                    <Metrics />
                    <Timer />
                </div>
                <div className="game-info-container">
                    {/* content for the game info goes here */}
                </div>
            </div>
        </div>
    );
}
