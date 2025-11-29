package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	_ "github.com/mattn/go-sqlite3"
)

type Contact struct {
	Name   string
	Number string
}

func main() {
	// CLI flags
	numbersPath := flag.String("n", "", "Path to numbers CSV file")
	messagePath := flag.String("m", "", "Path to message template file")
	fullMode := flag.Bool("full", false, "Send to all contacts (default: test mode - first contact only)")
	delayRange := flag.String("delay", "15-35", "Delay range in seconds (e.g., 15-35)")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	// Setup default directory in Downloads
	homeDir, _ := os.UserHomeDir()
	defaultDir := filepath.Join(homeDir, "Downloads", "whatsapp-broadcast")
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create directory: %v\n", err)
		os.Exit(1)
	}
	
	numbersFilePath := filepath.Join(defaultDir, "numbers.csv")
	messageFilePath := filepath.Join(defaultDir, "message.txt")
	
	// Check if sample files need to be created
	needSetup := false
	
	// Copy numbers.csv if it doesn't exist
	if _, err := os.Stat(numbersFilePath); os.IsNotExist(err) {
		needSetup = true
		fmt.Println("🔧 First time setup detected!")
		numbersContent := []byte("name,number\nJohn Doe,+919705937595\n")
		if err := os.WriteFile(numbersFilePath, numbersContent, 0644); err != nil {
			fmt.Printf("❌ Failed to create numbers.csv: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Created: %s\n", numbersFilePath)
	}
	
	// Copy message.txt if it doesn't exist
	if _, err := os.Stat(messageFilePath); os.IsNotExist(err) {
		if !needSetup {
			fmt.Println("🔧 First time setup detected!")
			needSetup = true
		}
		messageContent := []byte("Namaskaram ${name} 🙂\n\n🌟 This is a test. ✨  \n\nWith Joy,  \nIsha Volunteers 🙏\n")
		if err := os.WriteFile(messageFilePath, messageContent, 0644); err != nil {
			fmt.Printf("❌ Failed to create message.txt: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Created: %s\n", messageFilePath)
	}
	
	// Pause for user to customize files if setup was needed
	if needSetup {
		fmt.Println()
		fmt.Println("📝 Please customize the files with your contacts and message.")
		fmt.Print("   Press Enter to continue...")
		scanner.Scan()
		fmt.Println()
	}

	// Set default paths based on chosen directory
	if *numbersPath == "" {
		*numbersPath = numbersFilePath
	}
	if *messagePath == "" {
		*messagePath = messageFilePath
	}

	fmt.Printf("📋 Using numbers file: %s\n", *numbersPath)
	fmt.Printf("💬 Using message file: %s\n\n", *messagePath)

	// Check if files exist
	if _, err := os.Stat(*numbersPath); os.IsNotExist(err) {
		fmt.Printf("❌ Numbers file not found: %s\n", *numbersPath)
		os.Exit(1)
	}

	if _, err := os.Stat(*messagePath); os.IsNotExist(err) {
		fmt.Printf("❌ Message file not found: %s\n", *messagePath)
		os.Exit(1)
	}

	// Prompt for mode if not provided via flag
	if !*fullMode && !isFlagPassed("full") {
		fmt.Print("🎯 Send to all contacts? (y/N): ")
		scanner.Scan()
		input := strings.TrimSpace(strings.ToLower(scanner.Text()))
		*fullMode = (input == "y" || input == "yes")
	}

	fmt.Println() // Empty line for readability

	// Parse delay range
	var minDelay, maxDelay int
	_, err := fmt.Sscanf(*delayRange, "%d-%d", &minDelay, &maxDelay)
	if err != nil || minDelay <= 0 || maxDelay <= 0 || minDelay > maxDelay {
		fmt.Printf("❌ Invalid delay format. Use format: 15-35\n")
		os.Exit(1)
	}

	// Read message template
	messageBytes, err := os.ReadFile(*messagePath)
	if err != nil {
		fmt.Printf("❌ Failed to read message file: %v\n", err)
		os.Exit(1)
	}
	messageTemplate := strings.TrimSpace(string(messageBytes))

	// Read contacts
	contacts, err := readContacts(*numbersPath)
	if err != nil {
		fmt.Printf("❌ Failed to read contacts: %v\n", err)
		os.Exit(1)
	}

	if len(contacts) == 0 {
		fmt.Printf("❌ No valid contacts found in CSV file\n")
		os.Exit(1)
	}

	// Test mode or full mode
	testContacts := contacts
	modeText := "FULL MODE - all contacts"
	if !*fullMode {
		if len(contacts) > 0 {
			testContacts = contacts[:1]
		}
		modeText = "TEST MODE - first contact only"
	}

	fmt.Printf("📋 Found %d contacts, running in %s\n", len(contacts), modeText)
	fmt.Printf("⏱️  Delay between messages: %d-%d seconds\n\n", minDelay, maxDelay)
	
	// Truncate message for display
	displayMsg := messageTemplate
	if len(displayMsg) > 100 {
		displayMsg = displayMsg[:100] + "..."
	}
	fmt.Printf("Message template:\n\"%s\"\n\n", displayMsg)

	// Setup WhatsApp client
	client, err := setupWhatsAppClient()
	if err != nil {
		fmt.Printf("❌ Failed to setup WhatsApp client: %v\n", err)
		os.Exit(1)
	}
	defer client.Disconnect()

	// Wait for connection to be fully ready
	fmt.Println("⏳ Waiting for WhatsApp connection to stabilize...")
	time.Sleep(3 * time.Second)

	if !*fullMode {
		fmt.Printf("🚀 Starting test broadcast...\n\n")
	} else {
		fmt.Printf("🚀 Starting full broadcast...\n\n")
	}

	// Send messages
	successCount := 0
	failCount := 0

	for i, contact := range testContacts {
		personalizedMessage := strings.ReplaceAll(messageTemplate, "${name}", contact.Name)
		
		fmt.Printf("   Attempting to send to: %s\n", cleanNumber(contact.Number))
		
		err := sendMessage(client, contact.Number, personalizedMessage)
		if err != nil {
			fmt.Printf("❌ [%d/%d] Failed for %s (%s)\n", i+1, len(testContacts), contact.Name, contact.Number)
			fmt.Printf("   Error: %v\n", err)
			failCount++
		} else {
			fmt.Printf("✅ [%d/%d] Sent to %s (%s)\n", i+1, len(testContacts), contact.Name, contact.Number)
			successCount++
		}

		// Random delay between messages (except for last one)
		if i < len(testContacts)-1 {
			delay := time.Duration(minDelay+rand.Intn(maxDelay-minDelay+1)) * time.Second
			fmt.Printf("   ⏳ Waiting %ds before next message...\n\n", int(delay.Seconds()))
			time.Sleep(delay)
		}
	}

	if !*fullMode {
		fmt.Printf("\n🎉 Test broadcast completed!\n")
	} else {
		fmt.Printf("\n🎉 Full broadcast completed!\n")
	}
	fmt.Printf("📊 Summary: %d successful, %d failed\n", successCount, failCount)
}

// Helper function to check if a flag was passed
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func readContacts(filePath string) ([]Contact, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var contacts []Contact
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}
		if len(record) >= 2 {
			name := strings.TrimSpace(record[0])
			number := strings.TrimSpace(record[1])
			if name != "" && number != "" {
				contacts = append(contacts, Contact{
					Name:   name,
					Number: number,
				})
			}
		}
	}

	return contacts, nil
}

func setupWhatsAppClient() (*whatsmeow.Client, error) {
	// Create data directory for session storage
	dataDir := "./whatsapp-data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Setup database for storing session
	dbPath := filepath.Join(dataDir, "session.db")
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New(context.Background(), "sqlite3", "file:"+dbPath+"?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Get first device from store or create new
	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	clientLog := waLog.Stdout("Client", "ERROR", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	
	// Add event handler to track connection status
	connected := make(chan bool, 1)
	client.AddEventHandler(func(evt interface{}) {
		switch evt.(type) {
		case *events.Connected:
			connected <- true
		}
	})

	// If not logged in, show QR code
	if client.Store.ID == nil {
		fmt.Println("\n⚠️  Session not found or expired. Please authenticate:")
		fmt.Println("📱 Scan this QR code in WhatsApp (Menu → Linked Devices):")
		
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("✅ Authenticated successfully (session will be saved)")
				break
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}
		fmt.Println("✅ WhatsApp client ready (using saved session)")
	}
	
	// Wait for connection to be established
	select {
	case <-connected:
		fmt.Println("✅ Connection established")
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("connection timeout")
	}

	return client, nil
}

func cleanNumber(number string) string {
	cleanNum := strings.ReplaceAll(number, " ", "")
	cleanNum = strings.ReplaceAll(cleanNum, "-", "")
	cleanNum = strings.ReplaceAll(cleanNum, "+", "")
	cleanNum = strings.ReplaceAll(cleanNum, "(", "")
	cleanNum = strings.ReplaceAll(cleanNum, ")", "")
	return cleanNum
}

func sendMessage(client *whatsmeow.Client, number, message string) error {
	// Clean and format number
	cleanNum := cleanNumber(number)

	// Create JID (WhatsApp ID)
	jid := types.NewJID(cleanNum, types.DefaultUserServer)

	// Check if number is registered on WhatsApp
	isOnWhatsApp, err := client.IsOnWhatsApp(context.Background(), []string{jid.String()})
	if err != nil {
		return fmt.Errorf("failed to verify number: %w", err)
	}
	
	if len(isOnWhatsApp) == 0 || !isOnWhatsApp[0].IsIn {
		return fmt.Errorf("number is not registered on WhatsApp")
	}

	// Use the verified JID from WhatsApp
	verifiedJID := isOnWhatsApp[0].JID

	// Create message
	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	// Send message
	_, err = client.SendMessage(context.Background(), verifiedJID, msg)
	return err
}
