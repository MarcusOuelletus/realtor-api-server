package global

import (
	"os"
)

// WorkingDirectory - path to project root folder (needed for index.html)
var WorkingDirectory = os.Getenv("RETS_SERVER_PATH")

// DATABASE_NAME is the name of the project's Mongo Database
const DATABASE_NAME = "DATABASE_NAME"

// DatabaseIP - The IP of the mongo VM
const DatabaseIP = "DATABASE_IP"

// QUERY_LIMIT is the number of properties returned per page of results
const QUERY_LIMIT = 36

// AREA_LIMIT is the number of areas/municipalities/communities returned per dropdown page
const AREA_LIMIT = 10
