/*
 * @copyright Copyright (C) 2023 José Eduardo Hernández Rodríguez
 * @license MIT License
 * All rights reserved.
 * This source code is licensed under the MIT License found in the
 * LICENSE file in the root directory of this source tree
 *
 * Author: José Eduardo Hernández Rodríguez
 * Email: eduher.rodp@gmail.com
 * Date: April 23, 2023
 *
 * Description: This file contains the implementation of a factory for creating
 * Firestore clients. along with a main function that retrieves documents from
 * the "log" collection using the singleton instance of firestore.Client.
 */

package main

import (
	"example.com/m/v2/httpserver"
	"log"
)

func main() {
	port := ":8080"
	err := httpserver.StartServer(port)
	if err != nil {
		log.Fatal("Error starting server: ", err.Error())
	}
}
