package main

import (
	"UnityServer/game"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"
)

type BroadcastMessage struct {
	Type     string    `json:"type"`
	Health   int       `json:"health"`
	Attack   int       `json:"attack"`
	Position game.Vec2 `json:"position"`
	UUID     int       `json:"uuid"`
	Username string    `json:"username,omitempty"`
}

type PlayerData struct {
	UUID     int       `json:"uuid"`
	Name     string    `json:"name"`
	Position game.Vec2 `json:"position"` // Assuming Vector2 is a type you have defined
	// Add any other fields you want to send
}

var players []*game.Player
var boss *game.Boss
var playerRockets []*game.PlayerRocket
var bossRockets []*game.BossRocket
var mu sync.Mutex
var clientConnections []net.Conn // 假设有一个包含客户端连接的列表
func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	// defer print server stopped
	defer fmt.Println("Server stopped")

	// Create a ticker that triggers every 200ms
	ticker := time.NewTicker(200 * time.Millisecond)
	go func() {
		for range ticker.C {
			updateGameState()
			broadcastGameState()
		}
	}()

	// print server is running
	fmt.Println("Server is running...")
	initializeGame()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			return
		}

		// Add the new connection to the clientConnections slice
		mu.Lock()
		clientConnections = append(clientConnections, conn)
		mu.Unlock()

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	mu.Lock()
	if boss == nil {
		resetGame()
	}
	player := game.NewPlayer("PlayerName")
	players = append(players, player)
	mu.Unlock()

	// Serialize the player data to JSON
	playerData := PlayerData{
		UUID:     player.UUID,
		Name:     player.Username,
		Position: player.Position,
		// Add any other fields you want to send
	}
	jsonData, err := json.Marshal(playerData)
	if err != nil {
		fmt.Println("Error serializing player data:", err)
		return
	}

	// Send the JSON data to the client
	conn.Write(jsonData)
	conn.Write([]byte("\n")) // End with a newline character

	fmt.Println("Player connected. Current number of players:", len(players))

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Connection closed by client")
			break
		} else if err != nil {
			fmt.Println("Error reading client request:", err)
			break
		}
		//fmt.Println("Raw line from client:", line)
		var req game.ClientRequest
		err = json.Unmarshal([]byte(line), &req)
		if err != nil {
			fmt.Println("Error decoding client request:", err)
			continue // skip to the next iteration
		}

		// Print req string after decoding
		//fmt.Printf("req string: %+v\n", req)

		switch req.Type {
		case "move":
			handleMove(req, player)
		case "attack":
			handleAttack(req, player)
		case "hit":
			handleHit(req, player)
		default:
			fmt.Println("Unknown request type:", req.Type)
		}
	}

	// Remove the connection from the clientConnections slice when done
	mu.Lock()
	for i, clientConn := range clientConnections {
		if clientConn == conn {
			clientConnections = append(clientConnections[:i], clientConnections[i+1:]...)
			break
		}
	}
	// Also remove the player from the players slice
	for i, p := range players {
		if p == player {
			players = append(players[:i], players[i+1:]...)
			fmt.Println("Player disconnected. Current number of players:", len(players))
			break
		}
	}
	mu.Unlock()
}

func handleMove(req game.ClientRequest, player *game.Player) {
	scale := float32(0.2)

	// Update nextMove based on the speed given by req.X and req.Y, scaled by 0.2
	player.NextMove.X = float32(req.X) * scale
	player.NextMove.Y = float32(req.Y) * scale

	//fmt.Printf("Player %d has position (%f, %f)\n", req.UUID, player.Position.X, player.Position.Y)
}

func handleAttack(req game.ClientRequest, player *game.Player) {
	fmt.Printf("Player %d attacked\n", req.UUID)
	// increase the player's NextAttack by 1
	player.NextAttack = 1
}

func handleHit(req game.ClientRequest, player *game.Player) {
	fmt.Printf("Player %d hit target %d\n", req.UUID, req.TargetUUID)
	// 这里可以处理玩家攻击其他玩家或被其他玩家攻击的逻辑
}

func initializeGame() {
	fmt.Println("Initializing the game...")
	boss = nil
	players = make([]*game.Player, 0)
	playerRockets = make([]*game.PlayerRocket, 0)
	bossRockets = make([]*game.BossRocket, 0)
}

func resetGame() {
	fmt.Println("Resetting the game...")
	// new boss
	boss = game.NewBoss()

	// set all players' health to 100
	for _, player := range players {
		player.Health = 100
	}
	fmt.Println("Game reset with a new boss")
}

func updateGameState() {
	mu.Lock()
	defer mu.Unlock()

	minX := float32(-10.0)
	maxX := float32(10.0)
	minY := float32(-4.5)
	maxY := float32(4.5)

	// Helper function to clamp values
	clamp := func(val, min, max float32) float32 {
		if val < min {
			return min
		}
		if val > max {
			return max
		}
		return val
	}
	// if boss is not nil and boss's health is less than 0, set boss to nil
	if boss != nil && boss.Health <= 0 {
		boss = nil
		//initializeGame()
	}

	// 更新boss的位置
	if boss != nil {

		boss.Update()
		boss.Position.X = clamp(boss.Position.X, minX, maxX)
		boss.Position.Y = clamp(boss.Position.Y, minY, maxY)
		boss.NextMove.X = 0
		boss.NextMove.Y = 0

		// shoot a rocket if NextAttack is equal to 1
		if boss.NextAttack >= 1 {
			// players length is guaranteed to be > 0
			if len(players) > 0 {
				// choose a random player to target
				target := players[rand.Intn(len(players))]
				// create a new rocket
				rocket := game.NewBossRocket(boss.Position, boss.UUID, target.UUID)
				// add the rocket to the bossRockets slice
				bossRockets = append(bossRockets, rocket)
			}
			// reset the boss's NextAttack to 0
			boss.NextAttack = 0
		}
	}

	// remove players with health <= 0
	filteredPlayers := make([]*game.Player, 0)
	for _, player := range players {
		if player != nil && player.Health > 0 {
			filteredPlayers = append(filteredPlayers, player)
		}
	}
	players = filteredPlayers

	// 更新每个玩家的位置
	for _, player := range players {
		if player != nil {
			player.Position.X = clamp(player.Position.X+player.NextMove.X, minX, maxX)
			player.Position.Y = clamp(player.Position.Y+player.NextMove.Y, minY, maxY)
			player.NextMove.X = 0
			player.NextMove.Y = 0

			//update player's attack
			if player.NextAttack > 0 {
				// create a new rocket
				rocket := game.NewPlayerRocket(player.Position, player.UUID)
				// add the rocket to the playerRockets slice
				playerRockets = append(playerRockets, rocket)
				// reset the player's NextAttack to 0
				player.NextAttack = 0
			}
		}
	}
	// remove rockets with health <= 0
	filteredPlayerRockets := make([]*game.PlayerRocket, 0)
	for _, rocket := range playerRockets {
		if rocket != nil && rocket.Health > 0 {
			filteredPlayerRockets = append(filteredPlayerRockets, rocket)
		}
	}
	playerRockets = filteredPlayerRockets
	// 更新玩家火箭的位置
	// if boss is not nil
	if boss != nil {
		for _, rocket := range playerRockets {
			if rocket != nil {
				// print the rocket's position
				//fmt.Printf("Player rocket %d has position (%f, %f)\n", rocket.UUID, rocket.Position.X, rocket.Position.Y)
				// move to the direction of the boss
				var bossDir = rocket.Position.Direction(boss.Position)
				rocket.Position.X = rocket.Position.X + bossDir.X*0.5
				rocket.Position.Y = rocket.Position.Y + bossDir.Y*0.5
				//rocket.Position.X = clamp(rocket.Position.X, minX, maxX)
				//rocket.Position.Y = clamp(rocket.Position.Y, minY, maxY)
				rocket.NextMove.X = 0
				rocket.NextMove.Y = 0

				// if the rocket is out of the screen, set health to 0
				if rocket.Position.X < minX || rocket.Position.X > maxX || rocket.Position.Y < minY || rocket.Position.Y > maxY {
					rocket.Health = 0
				}

				// check if the rocket hits the boss, range is 0.5
				if boss != nil && rocket.Position.Distance(boss.Position) < 0.5 {
					// set the rocket's health to 0
					rocket.Health = 0
					// decrease the boss's health by rocket's attack
					boss.Health = boss.Health - rocket.Attack
					// print the boss's health
					fmt.Printf("Boss health: %d\n", boss.Health)
					// if the boss's health is 0, remove the boss
					if boss.Health == 0 {
						//boss = nil
					}
				}
			}
		}
	} else {
		// set all player rockets' health to 0
		for _, rocket := range playerRockets {
			if rocket != nil {
				rocket.Health = 0
			}
		}
	}

	// remove rockets with health <= 0
	filteredBossRockets := make([]*game.BossRocket, 0)
	for _, rocket := range bossRockets {
		if rocket != nil && rocket.Health > 0 {
			filteredBossRockets = append(filteredBossRockets, rocket)
		}
	}
	bossRockets = filteredBossRockets
	// 更新boss火箭的位置
	// if boss is not nil
	if boss != nil {
		for _, rocket := range bossRockets {
			if rocket != nil {
				// try to get the player by UUID
				var player *game.Player
				for _, p := range players {
					if p.UUID == rocket.Target {
						player = p
						break
					}
				}

				// if the player is not found, set the rocket's health to 0
				if player == nil {
					rocket.Health = 0
					continue
				}

				// move to the direction of the player
				var playerDir = rocket.Position.Direction(player.Position)
				rocket.Position.X = rocket.Position.X + playerDir.X*0.2
				rocket.Position.Y = rocket.Position.Y + playerDir.Y*0.2

				// if the rocket is out of the screen, set health to 0
				if rocket.Position.X < minX || rocket.Position.X > maxX || rocket.Position.Y < minY || rocket.Position.Y > maxY {
					rocket.Health = 0
				}

				//rocket.Position.X = clamp(rocket.Position.X, minX, maxX)
				//rocket.Position.Y = clamp(rocket.Position.Y, minY, maxY)
				rocket.NextMove.X = 0
				rocket.NextMove.Y = 0

				// check if the rocket hits the player, range is 1.5
				if rocket.Position.Distance(player.Position) < 1.5 {
					// set the rocket's health to 0
					rocket.Health = 0
					// decrease the player's health by rocket's attack
					player.Health = player.Health - rocket.Attack
					// print the player's health
					fmt.Printf("Player %s health: %d\n", player.UUID, player.Health)
					// if the player's health is 0,
					if player.Health == 0 {
						//for i, p := range players {
						//	if p == player {
						//		players = append(players[:i], players[i+1:]...)
						//		break
						//	}
						//}
					}
				}
			}
		}
	} else {
		// set all boss rockets' health to 0
		for _, rocket := range bossRockets {
			if rocket != nil {
				rocket.Health = 0
			}
		}
	}
}

func broadcastGameState() {
	mu.Lock()
	defer mu.Unlock()

	messages := []BroadcastMessage{}

	// 添加Boss的信息
	if boss != nil {
		messages = append(messages, BroadcastMessage{
			Type:     "boss",
			Health:   boss.Health,
			Attack:   boss.Attack,
			Position: boss.Position,
			UUID:     boss.UUID,
		})
	}

	// 添加玩家的信息
	for _, player := range players {
		if player != nil {
			messages = append(messages, BroadcastMessage{
				Type:     "player",
				Health:   player.Health,
				Attack:   player.Attack,
				Position: player.Position,
				UUID:     player.UUID,
				Username: player.Username,
			})
		}
	}

	// 添加玩家火箭的信息
	for _, rocket := range playerRockets {
		if rocket != nil {
			messages = append(messages, BroadcastMessage{
				Type:     "playerRocket",
				Health:   rocket.Health,
				Attack:   rocket.Attack,
				Position: rocket.Position,
				UUID:     rocket.UUID,
			})
		}
	}

	// 添加boss火箭的信息
	for _, rocket := range bossRockets {
		if rocket != nil {
			messages = append(messages, BroadcastMessage{
				Type:     "bossRocket",
				Health:   rocket.Health,
				Attack:   rocket.Attack,
				Position: rocket.Position,
				UUID:     rocket.UUID,
			})
		}
	}

	// print the messages
	//fmt.Println(messages)

	// 将信息转换为JSON
	jsonData, err := json.Marshal(messages)
	if err != nil {
		fmt.Println("Error marshaling broadcast data:", err)
		return
	}

	// 向所有连接的客户端发送数据
	for _, clientConn := range clientConnections {
		_, err := clientConn.Write(jsonData)
		if err != nil {
			fmt.Println("Error sending data to client:", err)
		}
	}
}
