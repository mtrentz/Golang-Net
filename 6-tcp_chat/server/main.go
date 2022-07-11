package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type Room struct {
	name     string
	cond     *sync.Cond
	messages []string
	n        int // amount of msgs
}

func (room *Room) insertMessage(msg string) {
	room.cond.L.Lock()
	room.messages = append(room.messages, msg)
	room.n += 1
	room.cond.Broadcast()
	room.cond.L.Unlock()
}

func newRoom(name string) *Room {
	return &Room{
		name: name,
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

type ChatRooms struct {
	rooms map[string]*Room
}

func newChatRooms() *ChatRooms {
	return &ChatRooms{
		rooms: make(map[string]*Room),
	}
}

// Try to create the room, if sucessfull, return the room.
func (cr *ChatRooms) createRoom(name string) (*Room, error) {
	// Check if already exists
	if _, ok := cr.rooms[name]; ok {
		fmt.Println("Room", name, "already exists")
		return nil, fmt.Errorf("Room %s already exists", name)
	}

	// Check if name not empty
	if name == "" {
		return nil, fmt.Errorf("Room name cannot be empty")
	}

	r := newRoom(name)

	// Create the room
	cr.rooms[name] = r

	return cr.rooms[name], nil
}

// Try to find the room, if sucessfull, return the room.
func (cr *ChatRooms) findRoom(name string) (*Room, error) {
	// Check if already exists
	if _, ok := cr.rooms[name]; !ok {
		fmt.Println("Room", name, "does not exist")
		return nil, fmt.Errorf("Room %s does not exist", name)
	}

	return cr.rooms[name], nil
}

// Return a list of all the rooms and the amount of msgs in each room.
func (cr *ChatRooms) listRooms() [][]string {
	var roomInfo [][]string
	for name, room := range cr.rooms {
		roomInfo = append(roomInfo, []string{name, fmt.Sprintf("%d", room.n)})
	}
	return roomInfo
}

// Check if the msg is a valid comment.
// This is done by checking if the msg starts with a '/'
// and then the next word is a valid 'command', defined on main
func isCommand(msg string) bool {
	if msg[0] != '/' {
		return false
	}

	// Removes the / from the msg, then separate by space to check the
	// first word, which supposed to be the command
	firstParam := strings.Fields(msg[1:])[0]
	for _, cmd := range commands {
		if firstParam == cmd {
			return true
		}
	}

	return false
}

var cr *ChatRooms
var commands []string

func main() {
	cr = newChatRooms()

	commands = []string{"create", "join", "list", "leave"}
	// _, _ = cr.createRoom("room1")

	l, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

// To properly send message from the client, I'll be using a pattern of
// string to separate username and the message.
func interpretMessage(msg string) (username string, message string) {
	// Split the message on *|*
	msgSplit := strings.Split(msg, "*|*")
	username = msgSplit[0]
	message = msgSplit[1]
	return username, message
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// First instantiate an empty room
	room := &Room{}

	// Store 'n', which is the amount of msgs
	// on the connected room
	var n int

	// Start goroutine to listen for user msgs
	// and commands to join/create room.
	go func() {
		// Keeps reading from the connection
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			msgReceived := scanner.Text()
			username, msg := interpretMessage(msgReceived)

			// Check if the msg is a command
			if isCommand(msg) {
				// Get the command/params, removing the first slash
				cmds := strings.Fields(msg[1:])

				if cmds[0] == "create" {
					// Check if has a next argument
					if len(cmds) <= 1 {
						io.WriteString(conn, "No room name provided."+"\n")
						continue
					}
					// Create the room, and assign the room to this user
					var err error
					room, err = cr.createRoom(cmds[1])
					if err != nil {
						io.WriteString(conn, fmt.Sprintf("Error creating room: %s\n", err)+"\n")
						// Set room to an empty room
						room = &Room{}
						continue
					}
				} else if cmds[0] == "join" {
					// Check if has a next argument
					if len(cmds) <= 1 {
						io.WriteString(conn, "No room name provided."+"\n")
						continue
					}
					// Join the room, and assign the room to this user
					var err error
					room, err = cr.findRoom(cmds[1])
					if err != nil {
						io.WriteString(conn, fmt.Sprintf("Error joining room: %s\n", err)+"\n")
						continue
					}
				} else if cmds[0] == "leave" {
					room = &Room{}
					io.WriteString(conn, "Left the room\n")
				} else if cmds[0] == "list" {
					// Listing all rooms will remove the user from the room
					// if he is in one. This is done because otherwise
					// any new msgs received on the client will override the
					// terminal and not show the list of rooms
					room = &Room{}

					roomInfo := cr.listRooms()

					roomInfoStr := ""
					// Send the list as "Name: name, Msgs: amount" separated by "|||"
					for _, info := range roomInfo {
						roomInfoStr += "Name: " + info[0] + ", Msgs: " + info[1] + "|||"
					}
					// Remove the last |||
					if roomInfoStr != "" {
						roomInfoStr = roomInfoStr[:len(roomInfoStr)-3]
					}

					// Check if there is a room
					if roomInfoStr == "" {
						io.WriteString(conn, "There are no rooms yet!"+"\n")
					} else {
						io.WriteString(conn, roomInfoStr+"\n")
					}
				} else {
					fmt.Println("There seems to be a command missing... Received:", msg)
				}
			} else if room.name != "" {
				// Empty room name means no room,
				// If the user has a room, then add the msg to it,
				// Here I'll insert with the username formatted as I wanted
				chatMsg := fmt.Sprintf("[%s]: %s", username, msg)
				room.insertMessage(chatMsg)
			} else {
				// If not in a room and not a command,
				// send msg to user to ask to join a room
				io.WriteString(conn, "You're not in a room, please join or create one.\n")
			}
		}
	}()

	// time.Sleep(10 * time.Second)

	// Now I start an inifnite loop that will
	// send all msgs in the chat room to the user
	// whenever there is a new msg.
	for {
		// If the user doesn't have a room, can't read msgs
		if room.name == "" {
			continue
		}

		// Wait for a new msg
		room.cond.L.Lock()
		for room.n == n {
			room.cond.Wait()
		}
		// Update the amount of msgs the user last saw from the room
		n = room.n

		// Send all messages to client as a single string.
		// Separate them by 3 pipes |
		msgs := strings.Join(room.messages, "|||")
		io.WriteString(conn, msgs+"\n")

		room.cond.L.Unlock()
	}
}
