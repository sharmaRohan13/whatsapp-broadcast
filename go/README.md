# WhatsApp Broadcast - Go Implementation

Cross-platform WhatsApp broadcast tool written in Go using `whatsmeow` library.

## Features

- ✅ Cross-platform (Windows, macOS, Linux)
- ✅ Single binary executable
- ✅ Persistent login (session saved in SQLite database)
- ✅ Test mode and full broadcast mode
- ✅ Customizable message delays
- ✅ Message personalization with ${name} placeholder

## Installation

### Prerequisites

- Go 1.21 or higher
- GCC compiler (for SQLite)
  - **macOS**: Xcode Command Line Tools (`xcode-select --install`)
  - **Windows**: MinGW-w64 or TDM-GCC
  - **Linux**: `sudo apt-get install build-essential` (Ubuntu/Debian)

### Install Dependencies

```bash
cd go
go mod download
go mod tidy
```

## Usage

### Quick Start (Recommended)

**First Run:**
1. Run the program - it will automatically create sample files in `~/Downloads/whatsapp-broadcast/`
2. Edit the files with your contacts and message
3. Press Enter to continue
4. Scan QR code when prompted
5. Messages will be sent in test mode (first contact only)

```bash
go run main.go
```

**Subsequent Runs:**
- The program will use existing files from `~/Downloads/whatsapp-broadcast/`
- No need to specify file paths
- Add actual message to be sent and the real contacts to send them to

### Run Modes

**Test Mode (sends to first contact only):**
```bash
go run main.go
```

**Full Mode (sends to all contacts):**
```bash
go run main.go -full
```

**Custom files:**
```bash
go run main.go -n path/to/numbers.csv -m path/to/message.txt
```

**Custom delay:**
```bash
go run main.go -delay 10-30
```

### Build Executable

**For your current OS:**
```bash
go build -o whatsapp-broadcast
```

**For Windows (from Mac/Linux):**
```bash
GOOS=windows GOARCH=amd64 go build -o whatsapp-broadcast.exe
```

**For macOS (from Windows/Linux):**
```bash
GOOS=darwin GOARCH=amd64 go build -o whatsapp-broadcast-mac
```

**For Linux (from Mac/Windows):**
```bash
GOOS=linux GOARCH=amd64 go build -o whatsapp-broadcast-linux
```

### Run the Executable

```bash
# macOS/Linux
./whatsapp-broadcast

# Windows
whatsapp-broadcast.exe
```

## Command Line Flags

```
-n string
    Path to numbers CSV file (default: ~/Downloads/whatsapp-broadcast/numbers.csv)
-m string
    Path to message template file (default: ~/Downloads/whatsapp-broadcast/message.txt)
-full
    Send to all contacts (default: test mode - first contact only)
-delay string
    Delay range in seconds (default: "15-35")
```

## Default File Location

On first run, the program automatically creates sample files in:
- **macOS/Linux**: `~/Downloads/whatsapp-broadcast/`
- **Windows**: `C:\Users\<username>\Downloads\whatsapp-broadcast\`

Files created:
- `numbers.csv` - Contact list with names and phone numbers
- `message.txt` - Message template with ${name} placeholder

## File Formats

### numbers.csv
```csv
name,number
John Doe,+1234567890
Jane Smith,+9876543210
```

### message.txt
```
Hello ${name}!

This is your personalized message.

Thanks!
```

## First-Time Setup

1. Run the application
2. Sample files will be created in `~/Downloads/whatsapp-broadcast/`
3. Edit the files with your contacts and message
4. Press Enter to continue
5. A QR code will appear in the terminal
6. Open WhatsApp on your phone
7. Go to: Menu → Linked Devices → Link a Device
8. Scan the QR code
9. Session is saved in `whatsapp-data/session.db`
10. Next runs won't require QR scanning or file setup

## Session Storage

- Session data is stored in `./whatsapp-data/session.db`
- The session persists across runs
- To reset: delete the `whatsapp-data` folder

## Examples

```bash
# First run - creates sample files and runs in test mode
go run main.go

# Full broadcast with existing files
go run main.go -full

# Full broadcast with custom delay
go run main.go -full -delay 20-40

# Custom files in test mode
go run main.go -n ~/my-contacts.csv -m ~/my-message.txt

# Build and run executable
go build -o broadcast
./broadcast -full
```

## Troubleshooting

**"gcc not found" or build errors:**
- Install GCC compiler for your platform
- macOS: `xcode-select --install`
- Windows: Install MinGW-w64
- Linux: `sudo apt-get install build-essential`

**"Failed to connect" errors:**
- Check your internet connection
- Delete `whatsapp-data` folder and re-authenticate
- Make sure WhatsApp is active on your phone

**Messages not sending:**
- Verify phone numbers are registered on WhatsApp
- Check number format (include country code)
- Ensure numbers are not blocked

## Advantages over Node.js Version

- ⚡ Faster execution
- 💾 Lower memory usage (~10-20MB vs ~100-200MB)
- 📦 Single executable file (no Node.js required)
- 🔒 Better session persistence
- 🚀 Easier deployment (just copy the binary)

## License

ISC
