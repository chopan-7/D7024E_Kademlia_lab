package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Message body is used to stora any information that we want to send in an RPC
type Msgbody struct {
	IP   string
	Data []byte // Hashed key value
}

type Response struct {
	RPC  string // String representing what kind of rpc the message is
	Body Msgbody
}

/*
	CLI module is detached from the rest of the service hence it
	has its own main package and function

	The CLI will only communicate outwards to the running kademlia service
	over port 10002 and not recieve any incoming messages.
*/
func main() {
	localIP := GetOutboundIP().String()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(`Welcome to hecnet!: 
Available commands:
ping <IP-address> 
put <file> 	(coming soon)
get <hash> 	(coming soon)
exit 		terminates program
`)
	for {
		fmt.Print("Enter command:")
		text, _ := reader.ReadString('\n') //read from terminal
		text = strings.TrimRight(text, "\n")

		words := strings.Fields(text)

		err := parser(localIP, words)

		if err != nil {
			fmt.Println(errors.Wrap(err, "Failed to parse command"))
		}
	}
}

func parser(localIP string, cmd []string) error {
	port := 10002
	if len(cmd) > 2 {
		return errors.New("Invalid command!")
	}
	switch command := cmd[0]; command {
	case "ping":
		ipAddress := []byte(cmd[1])
		ping(localIP, string(ipAddress), port)
	case "put":
		put(localIP, port, cmd[1])
	case "get":
		get(localIP, port, cmd[1])
	case "exit":
		exit(localIP, port)
	default:
		return errors.New(command + "is not a valid command...")
	}
	return nil
}

func ping(localIP, ip string, port int) error {
	ip_addr := []byte(ip)
	res := Response{
		RPC: "ping",
		Body: Msgbody{
			Data: ip_addr,
		},
	}

	marshalledReq := marshallData(res)

	addr := net.ParseIP(localIP)
	fmt.Println(addr)

	server := net.UDPAddr{
		Port: port,
		IP:   addr,
	}

	Conn, err := net.DialUDP("udp", nil, &server)

	if err != nil {
		return errors.Wrap(err, "CLI: Failed to open connection to local ip:"+localIP)
	}

	defer Conn.Close()

	Conn.Write([]byte(marshalledReq))
	buf := make([]byte, 1024)
	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallData(buf[0:n])

	fmt.Println("Received RPC: ", rec.RPC, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)

	return nil
}

func put(localIP string, port int, filePath string) error {
	dat, err := os.ReadFile(filePath)

	if err != nil {
		errors.Wrap(err, "Failed to read from file at:"+filePath)
	}
	h := sha1.New()
	h.Write(dat)
	hashFile := h.Sum(nil)

	req := Response{
		RPC: "put",
		Body: Msgbody{
			Data: hashFile,
		},
	}

	marshalledReq := marshallData(req)

	addr := net.ParseIP(localIP)

	server := net.UDPAddr{
		Port: port,
		IP:   addr,
	}

	Conn, err := net.DialUDP("udp", nil, &server)

	if err != nil {
		return errors.Wrap(err, "CLI: Failed to open connection to "+localIP)
	}

	defer Conn.Close()

	Conn.Write([]byte(marshalledReq))
	buf := make([]byte, 1024)
	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallData(buf[0:n])

	fmt.Println("Received RPC: ", rec.RPC, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)

	return nil
}

func get(localIP string, port int, data string) error {
	h := sha1.New()
	h.Write([]byte(data))
	hashFile := h.Sum(nil)

	res := Response{
		RPC: "get",
		Body: Msgbody{
			Data: hashFile,
		},
	}

	marshalledReq := marshallData(res)

	addr := net.ParseIP(localIP)

	server := net.UDPAddr{
		Port: port,
		IP:   addr,
	}

	Conn, err := net.DialUDP("udp", nil, &server)

	if err != nil {
		return errors.Wrap(err, "CLI: Failed to open connection to local ip:"+localIP)
	}

	defer Conn.Close()

	Conn.Write([]byte(marshalledReq))
	buf := make([]byte, 1024)
	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallData(buf[0:n])

	fmt.Println("Received RPC: ", rec.RPC, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)

	return nil
}

func exit(localIP string, port int) error {
	res := Response{
		RPC:  "exit",
		Body: Msgbody{},
	}

	marshalledReq := marshallData(res)

	addr := net.ParseIP(localIP)
	fmt.Println(addr)

	server := net.UDPAddr{
		Port: port,
		IP:   addr,
	}

	Conn, err := net.DialUDP("udp", nil, &server)

	if err != nil {
		return errors.Wrap(err, "CLI: Failed to open connection to local ip:"+localIP)
	}

	defer Conn.Close()

	Conn.Write([]byte(marshalledReq))
	buf := make([]byte, 1024)
	n, remoteaddr, _ := Conn.ReadFromUDP(buf)
	rec := unmarshallData(buf[0:n])

	fmt.Println("Received RPC: ", rec.RPC, "\nBody: ", rec.Body, "\nFrom ", remoteaddr)
	os.Exit(1)
	return nil
}

// Will marshall the response object into json
func marshallData(data Response) []byte {
	marshalledData, _ := json.Marshal(data)
	return marshalledData
}

func unmarshallData(data []byte) Response {
	var unmarshalledData Response
	json.Unmarshal([]byte(data), &unmarshalledData)
	return unmarshalledData
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
