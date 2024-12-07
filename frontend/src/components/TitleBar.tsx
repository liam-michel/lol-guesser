import l9 from "../../public/l9.png"
export default function TitleBar (){
    return(
        <>
        <div className="title-bar">
            <h1>League of Legends Guesser</h1>
            <img src={l9} alt="League of Legends" width={1000} height={500} />
        </div>
        </>
    )
}