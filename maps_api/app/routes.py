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
                        "name": columns[1],
                        "distance": distance
                    })
                # Print progress
                print(f"\r{idx + 1}/{len(os.listdir('data/routes'))}", end="")

    return nearest_stations


def get_nearest_station_info(location):
    nearest_station = {}
    nearest_distance = float('inf')
    nearest_route_name = ""

    with ThreadPoolExecutor() as executor:
        futures = []
        for filename in os.listdir("data/routes"):
            path = os.path.join("data/routes", filename)
            with open(path) as file:
                for line in file:
                    columns = line.split(",")
                    station_location = Location(float(columns[2]), float(columns[3]))
                    distance = distance_between_locations(location, station_location)
                    futures.append((distance, columns[1], filename))

        for future in futures:
            distance, station_name, route_name = future
            if distance < nearest_distance:
                nearest_distance = distance
                nearest_station = {
                    "name": station_name,
                    "distance": distance
                }
                nearest_route_name = route_name

    nearest_station["route_name"] = nearest_route_name
    return nearest_station


@app.route("/google-maps", methods=["GET"])
def get_routes():
    address = request.args.get("address")
    destination = request.args.get("destination")

    origin_location = geocode(address)
    destination_location = geocode(destination)

    if not origin_location or not destination_location:
        return jsonify({"error": "Invalid address or destination"}), 400

    route = calculate_route(origin_location, destination_location)

    nearest_station_info = get_nearest_station_info(origin_location)

    # Obtener la última etapa de la ruta
    last_step = route["routes"][0]["legs"][0]["steps"][-1]
    last_step_location = Location(last_step["end_location"]["lat"], last_step["end_location"]["lng"])

    # Obtener la parada más cercana al destino
    destination_station_info = get_nearest_station_info(destination_location)

    # Verificar si la ruta de origen y destino son iguales
    if nearest_station_info["route_name"] == destination_station_info["route_name"]:
        response = {
            "start_address": route["routes"][0]["legs"][0]["start_address"],
            "end_address": route["routes"][0]["legs"][0]["end_address"],
            "nearest_station_info": nearest_station_info,
            "destination_station_info": destination_station_info
        }
    else:
        response = {
            "start_address": route["routes"][0]["legs"][0]["start_address"],
            "end_address": route["routes"][0]["legs"][0]["end_address"],
            "nearest_station_info": nearest_station_info,
            "destination_station_info": {
                "distance": None,
                "name": "No hay una parada cercana en la misma ruta",
                "route_name": None
            }
        }
    mediumWebhook(response)
    return jsonify(response)

# Send the data to the medium webhook
def mediumWebhook(data):
    url = "http://localhost:3000/webhook"
    # Headers
    headers = {
        "Content-Type": "application/json",
        "X-Origin": "google-maps"
    }
    # Send the request
    response = requests.post(url, data=data)
