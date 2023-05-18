import os

from flask import request, jsonify
import requests
from app import app
# Import the GOOGLE_MAPS_API_KEY from config.py
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


# calculate_route function to get the best route between two locations using Google Maps Directions API
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


def distance_between_locations(location):
    # Initialize the shortest distance to a very large number
    shortest_distance = float('inf')
    nearest_station = None
    # Calculate the distance between two locations using Distance Matrix API from Google Maps
    # Construct the URL for the GET request to the Distance Matrix API
    distance_matrix_url = "https://maps.googleapis.com/maps/api/distancematrix/json"

    # Read all the files in the data/routes folder
    for filename in os.listdir("data/routes"):
        # Construct the path to open the file
        path = os.path.join("data/routes", filename)
        # Open the file
        with open(path) as file:
            # Read the file
            for line in file:
                # Split the line into a list of columns
                columns = line.split(",")
                # Get the latitude and longitude of the station
                station_location = Location(columns[2], columns[3])
                # Construct the parameters for the Distance Matrix API
                params = {
                    "origins": f"{location.latitude},{location.longitude}",
                    "destinations": f"{station_location.latitude},{station_location.longitude}",
                    "key": GOOGLE_MAPS_API_KEY
                }
                # Call the Distance Matrix API
                response = requests.get(distance_matrix_url, params=params)
                response_json = response.json()
                # Check if the response contains distance information
                if "rows" in response_json and len(response_json["rows"]) > 0:
                    elements = response_json["rows"][0].get("elements")
                    if elements and len(elements) > 0:
                        distance = elements[0].get("distance", {}).get("value")
                        # Compare the distance with the shortest distance
                        if distance and distance < shortest_distance:
                            shortest_distance = distance
                            nearest_station = {
                                "name": columns[1],
                                "latitude": columns[2],
                                "longitude": columns[3],
                                "route": filename,
                            }

    return nearest_station


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

        # Here we are calculating the nearest station from the origin
        # call function distance_between_locations with start_location as parameter
        # nearest_station_from_origin and nearest_station_from_destination must belong to the same route

        start_location = Location(response["start_location"]["latitude"], response["start_location"]["longitude"])
        end_location = Location(response["end_location"]["latitude"], response["end_location"]["longitude"])
        nearest_station_from_origin = distance_between_locations(start_location)
        nearest_station_from_destination = distance_between_locations(end_location)
        response["nearest_station_from_origin"] = nearest_station_from_origin
        response["nearest_station_from_destination"] = nearest_station_from_destination

        return jsonify(response)

    else:
        return jsonify({"error": "Invalid locations"})
