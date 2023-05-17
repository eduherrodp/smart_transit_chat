from flask import request, jsonify
import requests
from app import app
# Import the GOOGLE_MAPS_API_KEY from config.py.py
from config import GOOGLE_MAPS_API_KEY


# Location class to store latitude and longitude
class Location:
    def __init__(self, latitude, longitude):
        self.latitude = latitude
        self.longitude = longitude


# geocode function to get coordinates from an address using Google Maps Geocoding API
def geocode(address):
    # Construct the URL for the GET request to the Geocoding API
    geocoding_url = "https://maps.googleapis.com/maps/api/geocode/json"
    params = {
        "address": address,
        "key": GOOGLE_MAPS_API_KEY
    }

    # Call the Geocoding API
    response = requests.get(geocoding_url, params=params)
    response_json = response.json()

    if response_json["status"] == "OK" and len(response_json["results"]) > 0:
        # Extract the coordinates
        location_data = response_json["results"][0]["geometry"]["location"]
        location = Location(location_data["lat"], location_data["lng"])
        return location

    return None


# calculateRoute function to get the best route between two locations using Google Maps Directions API
def calculate_route(origin, destination):
    # Construct the URL for the GET request to the Directions API
    directions_url = "https://maps.googleapis.com/maps/api/directions/json"
    params = {
        "origin": f"{origin.latitude},{origin.longitude}",
        "destination": f"{destination.latitude},{destination.longitude}",
        "key": GOOGLE_MAPS_API_KEY
    }

    # Call the Directions API
    response = requests.get(directions_url, params=params)
    response_json = response.json()
    return response_json


@app.route('/calculate-route', methods=['POST'])
def handle_calculate_route():
    data = request.get_json()
    current_location = geocode(data["current_location"])
    destination_location = geocode(data["destination_location"])

    if current_location is not None and destination_location is not None:
        directions_response = calculate_route(current_location, destination_location)
        # Construct the response
        response = {
            "start_address": directions_response["routes"][0]["legs"][0]["start_address"],
            "start_location": {
                "latitude": directions_response["routes"][0]["legs"][0]["start_location"]["lat"],
                "longitude": directions_response["routes"][0]["legs"][0]["start_location"]["lng"]
            },
            "end_address": directions_response["routes"][0]["legs"][0]["end_address"],
            "end_location": {
                "latitude": directions_response["routes"][0]["legs"][0]["end_location"]["lat"],
                "longitude": directions_response["routes"][0]["legs"][0]["end_location"]["lng"]
            },
            "distance": directions_response["routes"][0]["legs"][0]["distance"]["text"]
        }

        return jsonify(response)

    else:
        return jsonify({"error": "Invalid locations"})
