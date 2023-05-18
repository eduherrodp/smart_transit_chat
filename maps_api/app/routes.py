import os
from concurrent.futures import ThreadPoolExecutor

import requests
from flask import request, jsonify
from app import app
from config import GOOGLE_MAPS_API_KEY


class Location:
    def __init__(self, latitude, longitude):
        self.latitude = latitude
        self.longitude = longitude


def geocode(address):
    geocoding_url = "https://maps.googleapis.com/maps/api/geocode/json"
    params = {
        "address": address,
        "key": GOOGLE_MAPS_API_KEY
    }
    response = requests.get(geocoding_url, params=params)
    response_json = response.json()

    if response_json["status"] == "OK" and len(response_json["results"]) > 0:
        location_data = response_json["results"][0]["geometry"]["location"]
        location = Location(location_data["lat"], location_data["lng"])
        return location

    return None


def calculate_route(origin, destination):
    directions_url = "https://maps.googleapis.com/maps/api/directions/json"
    params = {
        "origin": f"{origin.latitude},{origin.longitude}",
        "destination": f"{destination.latitude},{destination.longitude}",
        "key": GOOGLE_MAPS_API_KEY
    }
    response = requests.get(directions_url, params=params)
    response_json = response.json()
    return response_json


def distance_between_locations(location1, location2):
    distance_matrix_url = "https://maps.googleapis.com/maps/api/distancematrix/json"
    params = {
        "origins": f"{location1.latitude},{location1.longitude}",
        "destinations": f"{location2.latitude},{location2.longitude}",
        "key": GOOGLE_MAPS_API_KEY
    }
    response = requests.get(distance_matrix_url, params=params)
    response_json = response.json()
    if "rows" in response_json and len(response_json["rows"]) > 0:
        elements = response_json["rows"][0].get("elements")
        if elements and len(elements) > 0:
            distance = elements[0].get("distance", {}).get("value")
            return distance
    return float('inf')


def get_nearest_stations(location):
    nearest_stations = []

    with ThreadPoolExecutor() as executor:
        futures = []
        for filename in os.listdir("data/routes"):
            path = os.path.join("data/routes", filename)
            with open(path) as file:
                for line in file:
                    columns = line.split(",")
                    station_location = Location(float(columns[2]), float(columns[3]))
                    futures.append(executor.submit(distance_between_locations, location, station_location))

        distances = [future.result() for future in futures]

    for idx, filename in enumerate(os.listdir("data/routes")):
        path = os.path.join("data/routes", filename)
        with open(path) as file:
            for line in file:
                columns = line.split(",")
                station_location = Location(float(columns[2]), float(columns[3]))
                distance = distances.pop(0)
                if distance < float('inf'):
                    nearest_stations.append({
                        "route": filename[:-4],
                        "stop": columns[0],
                        "name": columns[1],
                        "distance": distance
                    })
                # Print progress
                print(f"\r{idx + 1}/{len(os.listdir('data/routes'))}", end="")

    return nearest_stations


@app.route("/routes", methods=["GET"])
def get_routes():
    address = request.args.get("address")
    destination = request.args.get("destination")

    origin_location = geocode(address)
    destination_location = geocode(destination)

    if not origin_location or not destination_location:
        return jsonify({"error": "Invalid address or destination"}), 400

    route = calculate_route(origin_location, destination_location)

    nearest_stations = get_nearest_stations(origin_location)

    response = {
        "route": route,
        "nearest_stations": nearest_stations
    }

    return jsonify(response)
