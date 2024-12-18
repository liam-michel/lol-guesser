# .env File configuration
The .env file needs to be located in the root directory (i.e. same level as this README) \
## .env Variables
1. RIOT_GAMES_API_KEY - this needs to be your own riot games API key, not necessary for the use of the application but useful if you want to use some test scripts such as **scrape_champs.py**
2. VITE_GOLANG_PORT - We have to preface this env variable with VITE_ So that Vite can properly load the variable given its location. This variable is used to tell the frontend which port the golang http server is running on.
3. VITE_REACT_PORT - This is used in the vite configuration to tell it which port to run on. If not set, the server will default to port 3000 

## Basic project structure
### /backend
This directory contains all of the golang code for the http server. This is where all of the processing happens,and I've made a point of completely isolating all of the processing to the backend so that the frontend is light
### /frontend 
This contains the React project with Vite. 
### /static
This is where all of the static files such as champion images and **champions.json** are contained.
### go.mod 
This contains all of the dependencies that you'll need for the project
