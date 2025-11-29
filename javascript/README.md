# WhatsApp Broadcast CLI

A cross-platform command-line tool for sending personalized WhatsApp broadcast messages to multiple contacts.

## Features

- ✅ Cross-platform (Windows, macOS, Linux)
- 📱 QR code authentication
- 🎯 Test mode (send to first contact only)
- 🚀 Full broadcast mode
- ⏱️ Customizable delay between messages
- 💬 Message personalization with name placeholders
- 🎨 Colored terminal output

## Prerequisites

- Node.js 16.0.0 or higher
- WhatsApp account

## Installation

### Option 1: Local Installation (Recommended)

```bash
# Clone or navigate to the project directory
cd whatsapp-broadcast

# Install dependencies
npm install
```

### Option 2: Global Installation

```bash
npm install -g .
```

After global installation, you can run `whatsapp-broadcast` from anywhere.

## Usage

### Basic Commands

**Test Mode (sends to first contact only):**
```bash
node send.js -n numbers.csv -m message.txt
```

**Full Mode (sends to all contacts):**
```bash
node send.js -n numbers.csv -m message.txt --full
```

**Custom delay between messages:**
```bash
node send.js -n numbers.csv -m message.txt --delay 10-30
```

### Windows-Specific Commands

```cmd
# Using npm script
npm start -- -n numbers.csv -m message.txt

# Direct execution
node send.js -n numbers.csv -m message.txt --full
```

### macOS/Linux-Specific Commands

```bash
# Using npm script
npm start -- -n numbers.csv -m message.txt

# Direct execution
./send.js -n numbers.csv -m message.txt --full
```

### Command Options

```
Options:
  -V, --version          output the version number
  -n, --numbers <path>   Path to CSV file with contacts (name,number)
  -m, --message <path>   Path to text file with message template
  -f, --full             Send to all contacts (default: test mode)
  -d, --delay <range>    Delay range in seconds (e.g., "15-35") (default: "15-35")
  -h, --help             display help for command
```

## File Formats

### numbers.csv

CSV file with two columns: `name` and `number`

```csv
name,number
John Doe,+1234567890
Jane Smith,+9876543210
```

### message.txt

Plain text file with your message. Use `${name}` as a placeholder for personalization:

```
Hello ${name}!

This is a personalized message just for you.

Thanks!
```

## First-Time Setup

1. **Run the application** - On first run, a QR code will appear in your terminal
2. **Open WhatsApp** on your phone
3. **Navigate to:** Menu → Linked Devices → Link a Device
4. **Scan the QR code** displayed in the terminal
5. **Authentication saved** - You won't need to scan again on subsequent runs

## Cross-Platform Compatibility

This app works identically on:
- ✅ Windows 10/11
- ✅ macOS (Intel & Apple Silicon)
- ✅ Linux (Ubuntu, Debian, Fedora, etc.)

### Platform-Specific Notes

**Windows:**
- Use Command Prompt, PowerShell, or Windows Terminal
- File paths can use backslashes: `.\numbers.csv`

**macOS/Linux:**
- Use Terminal or any bash/zsh shell
- File paths use forward slashes: `./numbers.csv`
- Make send.js executable: `chmod +x send.js`

## Examples

### Example 1: Test Run
```bash
node send.js -n contacts.csv -m announcement.txt
```
Sends message to the first contact only (safe testing).

### Example 2: Full Broadcast with Custom Delay
```bash
node send.js -n contacts.csv -m announcement.txt --full --delay 20-40
```
Sends to all contacts with 20-40 second random delays.

### Example 3: Using Absolute Paths (Windows)
```cmd
node send.js -n C:\Users\YourName\contacts.csv -m C:\Users\YourName\message.txt --full
```

### Example 4: Using Absolute Paths (macOS/Linux)
```bash
node send.js -n ~/Documents/contacts.csv -m ~/Documents/message.txt --full
```

## Tips

1. **Always test first** - Run without `--full` flag to send to first contact only
2. **Respect WhatsApp limits** - Use appropriate delays (15-35 seconds recommended)
3. **Personalize messages** - Use `${name}` placeholder for better engagement
4. **Keep session active** - Authentication is cached in `.wwebjs_auth` folder
5. **Backup your data** - Keep backups of your CSV and message files

## Troubleshooting

**"Authentication failed"**
- Delete `.wwebjs_auth` folder and re-authenticate

**"No valid contacts found"**
- Check CSV format (must have `name` and `number` columns)
- Ensure no empty rows

**"Module not found"**
- Run `npm install` to install dependencies

**QR code not showing (Windows)**
- Ensure terminal supports Unicode characters
- Try Windows Terminal instead of Command Prompt

## Security

- Authentication data is stored locally in `.wwebjs_auth`
- Never commit `.wwebjs_auth` to version control
- Keep your contact lists private

## License

ISC
