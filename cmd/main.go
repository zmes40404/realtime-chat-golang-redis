package main
// "package main" is a special reserved name of Go	
// main.go is execution entry point of the entire project
// Only package main can compile into an executable file (binary)

import (
	"chatroom-project/internal/server"
	"log"
)

// "func main()" is the start function of the main program and must be in package main
// "go build", "go run" will look for the main() function to start the program
func main() {	// Gin + WebSocket Server starting point
	router := server.SetupRouter()
	log.Println("🚀 Server running at http://localhost:8080")
	router.Run(":8080")
}