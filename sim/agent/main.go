package main

import (
	"encoding/binary"
	"github.com/goccy/go-json"
	"github.com/wowsims/wotlk/sim"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight/dps"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
)
import _ "net/http/pprof"

func main() {
	socket := "/tmp/sim-agent.sock"

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// remove socket if it exists
	_, err := os.Stat(socket)
	if err == nil {
		err = os.Remove(socket)
		if err != nil {
			log.Println("Error removing socket:", err.Error())
			os.Exit(1)
		}
	}
	l, err := net.Listen("unix", "/tmp/sim-agent.sock")
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer func() {
		err = l.Close()
		if err != nil {
			log.Println("Error closing listener:", err.Error())
			os.Exit(1)
		}
	}()

	log.Println("Listening on " + socket)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			continue
		}

		log.Println("Connection established from ", conn.RemoteAddr().String())

		go handleRequest(conn)
	}
}

type Response struct {
	Success bool
	Body    map[string]interface{}
}

type Request struct {
	Command string
	Body    map[string]interface{}
}

type Session struct {
	sim   *core.Simulation
	agent Agent
}

var agentFactoryMap = map[string]func(simAgent core.Agent) Agent{
	"DpsDeathknight": func(simAgent core.Agent) Agent {
		return &DpsDeathknightAgent{simAgent: simAgent.(*dps.DpsDeathknight)}
	},
}

func getAgentFactory(agent core.Agent) Agent {
	agentName := reflect.TypeOf(agent).Elem().Name()

	if agentFactory, ok := agentFactoryMap[agentName]; ok {
		return agentFactory(agent)
	}

	log.Println("No agent factory found for agent: ", agentName)
	return nil
}

func startSimHandler(request *Request, session *Session) Response {
	raidSimRequestJSON, err := json.Marshal(request.Body["RaidSimRequest"])

	if err != nil {
		log.Println("Error marshalling:", err.Error())
		return Response{Success: false}
	}

	raidSimRequestProto := &proto.RaidSimRequest{}
	err = protojson.Unmarshal(raidSimRequestJSON, raidSimRequestProto)

	if err != nil {
		log.Println("Error unmarshalling:", err.Error())
		return Response{Success: false}
	}

	sim.RegisterAll()
	session.sim = core.CreateSim(raidSimRequestProto)
	session.sim.Start()

	char := session.sim.GetPlayer()
	char.GetCharacter().ClearGCDAction()
	session.agent = getAgentFactory(char)

	if session.agent == nil {
		return Response{Success: false}
	}

	session.agent.Init(session)
	return Response{Success: true}
}

func getStateHandler(request *Request, session *Session) Response {
	if session.sim == nil {
		return Response{Success: false}
	}
	return session.agent.GetState(session)
}

func castHandler(request *Request, session *Session) Response {
	if session.sim == nil {
		return Response{Success: false}
	}

	body := request.Body
	spell := body["spell"].(string)

	return session.agent.Cast(spell, session)
}

func waitDurationHandler(request *Request, session *Session) Response {
	if session.sim == nil {
		return Response{Success: false}
	}

	duration := int(request.Body["duration"].(float64))
	return session.agent.Wait(duration, session)
}

var handlers = map[string]func(request *Request, session *Session) Response{
	"START_SIM_SESSION": startSimHandler,
	"GET_STATE":         getStateHandler,
	"CAST":              castHandler,
	"WAIT_DURATION":     waitDurationHandler,
}

func handleRequest(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection:", err.Error())
		}
	}()

	readBuffer := make([]byte, 1024*20)
	writeBuffer := make([]byte, 1024*5)
	session := &Session{}

	for {
		size, err := conn.Read(readBuffer)

		if err != nil {
			if err.Error() == "EOF" {
				log.Println("Connection closed by client")
			} else {
				log.Println("Error reading:", err.Error())
			}

			return
		}

		if size < 4 {
			log.Println("Wrong size, expected 4 bytes not ", size)
			return
		}

		// Read 4 bytes from the buffer and convert to an int
		// This is the length of the message
		msgLen := int(binary.LittleEndian.Uint32(readBuffer[:4]))

		currentSize := size
		// keep reading until we have the full message
		for currentSize < msgLen+4 {
			size, err = conn.Read(readBuffer[size:])
			currentSize += size

			if err != nil {
				log.Println("Error reading:", err.Error())
				return
			}
		}
		//if msgLen != size-4 {
		//	log.Println("Wrong size, expected ", msgLen, " bytes not ", size)
		//	return
		//}

		requestBody := &Request{}
		err = json.Unmarshal(readBuffer[4:msgLen+4], &requestBody)

		if err != nil {
			log.Println("Error unmarshalling:", err.Error())
			return
		}

		//log.Println("Message Received:", requestBody)
		responseBody := Response{Success: false}

		if handler, ok := handlers[requestBody.Command]; ok {
			responseBody = handler(requestBody, session)
		} else {
			log.Println("No handler for command:", requestBody.Command)
		}

		// convert the response to json
		response, err := json.Marshal(responseBody)
		if err != nil {
			log.Println("Error marshalling:", err.Error())
			return
		}

		// Write the length of the message to the buffer
		binary.LittleEndian.PutUint32(writeBuffer, uint32(len(response)))
		copy(writeBuffer[4:], response)

		_, err = conn.Write(writeBuffer[:len(response)+4])
		if err != nil {
			log.Println("Error writing:", err.Error())
			return
		}
	}
}
