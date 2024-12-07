import TitleBar from "./TitleBar";
import Timer from "./Timer";
import Score from "./Score";
import UserInput from "./UserInput";
import RandomChampion from "./RandomChampion";

export default function GameLoop() {
    return (
        <div className="game-loop">
            <TitleBar />
            <RandomChampion />
            <div className="stats">
                <UserInput />
                <Score />
                <Timer />
            </div>
        </div>
    );
}
