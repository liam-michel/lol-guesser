export type Champion = null | {
    name: string;
    imageURL: string;
    
}

export type ScoreSubmission = {
    authToken: string;
    username: string;
    correctGuesses: number;
    incorrectGuesses: number;
    accuracy: number;

}