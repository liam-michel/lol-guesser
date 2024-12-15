//file with functions to set and retrieve scores from the backend
import type { ScoreSubmission } from '../types';
import { isTokenExpired } from './user_api';

export async function setNewScore(submission: ScoreSubmission){
  try{
    //destructure the auth token
    const {authToken} = submission;
    //check if the token is expired
    return;

  }
}