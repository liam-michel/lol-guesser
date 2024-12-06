import json
import os
import requests
from dotenv import load_dotenv

def read_champs():
    with open("champions.json", "r") as file:
        data = json.load(file)
        return list(data["data"].keys())


def scrape_images(names):


    # Create images directory if it doesn't exist
    os.makedirs("./static/images", exist_ok=True)

    # Load API key from .env file
    load_dotenv()
    api_key = os.getenv("RIOT_GAMES_API_KEY")

    for name in names:
        # Get champion data from Data Dragon API (no auth needed)
        url = f"http://ddragon.leagueoflegends.com/cdn/13.24.1/img/champion/{name}.png"
        
        response = requests.get(url)
        if response.status_code == 200:
            # Save image to static/images directory
            with open(f"./static/images/{name}.png", "wb") as f:
                f.write(response.content)
            print(f"Downloaded {name}.png")
        else:
            print(f"Failed to download {name}.png")

if __name__ == "__main__":
    names = read_champs()
    scrape_images(names)



