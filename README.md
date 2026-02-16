# learning-go
Learning go from sratch

# Go (Golang) Setup Guide – Windows

This guide covers:

- Installing Go on Windows
- Setting up VS Code
- Creating your first Go project
- Running a Hello World program
- Building an executable (.exe)

---

# 1️⃣ Install Go on Windows

## Step 1: Download Go

Go to the official website:

https://go.dev/dl/

Download:
> **Windows (.msi installer)**

---

## Step 2: Install

1. Run the downloaded `.msi` file  
2. Click **Next**
3. Accept license
4. Keep default install location
5. Click **Install**
6. Click **Finish**

---

## Step 3: Verify Installation

Open Command Prompt (CMD) or PowerShell:

```bash
go version
```

Expected output:

```bash
go version go1.xx.x windows/amd64
```

If this works, Go is installed correctly ✅

---

# 2️⃣ Install VS Code

Download from:

https://code.visualstudio.com/

Install with default settings.

---

# 3️⃣ Install Go Extension in VS Code

1. Open VS Code
2. Go to Extensions (`Ctrl + Shift + X`)
3. Search: **Go**
4. Install the extension by Google

If prompted, install recommended Go tools.

---

# 4️⃣ Create Your First Go Project

## Step 1: Create Project Folder

Example:

```
C:\go-projects\hello
```

---

## Step 2: Open Folder in VS Code

- Open VS Code
- Click **File → Open Folder**
- Select your project folder

---

# 5️⃣ Initialize Go Module

Open terminal in VS Code (`Ctrl + ~`)

Run:

```bash
go mod init hello
```

Output:

```bash
go: creating new go.mod: module hello
```

This creates a `go.mod` file.

---

# 6️⃣ Create Hello World Program

Create a file named:

```
main.go
```

Add this code:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

Save the file.

---

# 7️⃣ Run the Program

In terminal:

```bash
go run main.go
```

Output:

```
Hello, World!
```

✅ You have successfully run your first Go program.

---

# 8️⃣ Build Executable File (.exe)

To compile the program:

```bash
go build main.go
```

This creates:

```
main.exe
```

---

# 9️⃣ Run the Executable

Run from terminal:

```bash
.\main.exe
```

Or double-click `main.exe` in File Explorer.

---

# 🔁 Difference Between Commands

| Command | What It Does |
|----------|--------------|
| `go run main.go` | Compiles and runs immediately |
| `go build main.go` | Compiles and creates executable |
| `go version` | Shows installed Go version |
| `go mod init <name>` | Initializes a Go module |

---

# 🚀 Project Structure After Setup

```
hello/
│── go.mod
│── main.go
│── main.exe (after build)
```

---

# 🎉 You Are Ready!

You can now start learning:

- Variables
- Functions
- Structs
- Loops
- Goroutines
- Building APIs

Happy Coding with Go 🚀
