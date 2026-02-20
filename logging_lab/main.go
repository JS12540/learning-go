package main

import (
	"log"
	"os"
)

func main() {

	// 1️⃣ Configure logger
	log.SetPrefix("APP: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 2️⃣ Basic logs
	log.Println("Application started")

	// ------------------------------------------------

	// 3️⃣ Different "levels" (simulated)

	log.Println("[INFO] Server running")
	log.Println("[DEBUG] Debugging value x=42")
	log.Println("[WARN] Disk space low")
	log.Println("[ERROR] Database connection failed")

	// ------------------------------------------------

	// 4️⃣ Logging to file (production common)

	file, err := os.OpenFile("app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)

	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	log.SetOutput(file)

	log.Println("[INFO] This goes to file")
	log.Println("[ERROR] Something failed")
	log.Default().Println("[INFO] This goes to console")

	// ------------------------------------------------

	// 5️⃣ Fatal vs Panic

	// log.Fatal("Fatal error → exits program")

	// log.Panic("Panic error → crashes with stacktrace")
}
