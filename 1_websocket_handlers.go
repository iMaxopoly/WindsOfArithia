package main

import (
	"fmt"
	"encoding/hex"
	"strings"
	"time"
	"math/rand"
	"sync"

	// 3rd Party
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

// Partially borrowed code from https://github.com/domluna/websocket-golang-chat/blob/master/chat.go


// Map containing clients
var websocket__ActiveClients = make(map[ClientConn]int)
var websocket__ActiveClientsRWMutex sync.RWMutex

var websocket__ws_remote_connections []string
var websocket__ws_remote_connectionsRWMutex sync.RWMutex

// Client connection consists of the websocket and the client ip
type ClientConn struct {
	websocket *websocket.Conn
	clientIP  string
}

var upgrader = websocket.New(websocket__ws_handler)

// WS Route
func websocket__route(ctx *iris.Context) {
	if (string(ctx.UserAgent()) != conf__launcher_useragent) {
		util__debug_info(fmt.Sprintf("[WS]| WS-SERVER %s - Invalid U-A", ctx.Request.Header.Peek("X-Forwarded-For")))
		return
	}
	upgrader.Upgrade(ctx)
}

func websocket__ws_handler(ws *websocket.Conn) {
	defer func() {
		// Safety-check remove the IP from firewall
		chan_sentry <- struct_sentry{cmd:"Remove", ip: ws.Header("X-Forwarded-For"), rulename: "ArithiaGame"}

		// Disconnect the websocket connection
		if err := ws.Close(); err != nil {
			util__debug_info(fmt.Sprintf("[WS]| WS-SERVER Websocket could not be closed - %s", err.Error()))
		}
		for i, v := range websocket__ws_remote_connections {
			if(v == ws.Header("X-Forwarded-For")){
				websocket__ws_remote_connectionsRWMutex.Lock()
				websocket__ws_remote_connections = append(websocket__ws_remote_connections[:i], websocket__ws_remote_connections[i+1:]...)
				websocket__ws_remote_connectionsRWMutex.Unlock()
				break
			}
		}
	}()

	fmt.Println(ws.Headers())
	client := ws.Header("X-Forwarded-For")
	util__debug_info(fmt.Sprintf("[WS]| WS-SERVER Client connected: %s", client))
	sockCli := ClientConn{ws, client}
	websocket__ActiveClientsRWMutex.Lock()
	websocket__ActiveClients[sockCli] = 0
	websocket__ActiveClientsRWMutex.Unlock()
	util__debug_info(fmt.Sprintf("[WS]| WS-SERVER Number of clients connected - %d", len(websocket__ActiveClients)))

	for {
		mt, msg, w_err := ws.ReadMessage()
		if w_err != nil {
			// If we cannot Read then the connection is closed
			util__debug_info(fmt.Sprintf("[WS]| WS-SERVER Websocket Disconnected waiting - %s", w_err))
			// remove the ws client conn from our active clients
			websocket__delete_conn(sockCli, client)
			util__debug_info(fmt.Sprintf("[WS]| WS-SERVER Number of clients still connected - %d", len(websocket__ActiveClients)))
			return
		}

		// Check if this IP address is already connected once
		for _, v := range websocket__ws_remote_connections {
			if(v == client){
				util__debug_info(fmt.Sprintf("[WS]| WS-SERVER Already connected - %s", client))
				websocket__ActiveClientsRWMutex.Lock()
				delete(websocket__ActiveClients, sockCli)
				websocket__ActiveClientsRWMutex.Unlock()
				return
			}
		}

		// Recovering Username and Password
		// Converting from Hex to non-Hex?
		accesskey_xor, err := hex.DecodeString(string(msg))
		if err != nil {
			util__debug_info(fmt.Sprintf("[WS]| WS-SERVER %s - Weird AccessKey", client))
			websocket__delete_conn(sockCli, client)
			return
		}

		// Recovering from Xor'ed data to real data
		accesskey_plaintext := string(util__xor_encryptdecrypt(accesskey_xor))

		// Storing Username and Password to respective variables
		username, password, file_security_digest := func() (string, string, string) {
			result := strings.Split(accesskey_plaintext, "wind")
			return result[0], result[1], result[2]
		}()

		// Evaluating File security demo
		//if file_security_digest != server_config.files_digest{
		//	util__debug_info(fmt.Sprintf("[WS]| WS-SERVER %s - %s File Security check failed", client, file_security_digest))
		//	ws.WriteMessage(mt, []byte(conf__launcher_err_loginfail))
		//	websocket__delete_conn(sockCli, client)
		//	return
		//}

		// Evaluating user exists
		rPass, nEMID := database__get_password_and_nEMID(username)
		if rPass == "" {
			util__debug_info(fmt.Sprintf("[WS]| WS-SERVER %s - User doesn't exist", client))
			ws.WriteMessage(mt, []byte(conf__launcher_err_loginfail))
			websocket__delete_conn(sockCli, client)
			return
		}

		// Testing password
		if util__aes_decrypt([]byte(server_config.aes_key), rPass) != password {
			util__debug_info(fmt.Sprintf("[WS]| WS-SERVER %s - Wrong Password", client))
			ws.WriteMessage(mt, []byte(conf__launcher_err_loginfail))
			websocket__delete_conn(sockCli, client)
			return
		}

		// Generating token
		gamekey := (func(n int) (string) {
			letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			rand.Seed(time.Now().UTC().UnixNano())
			b := make([]rune, n)
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			return string(b)
		}(40) + util__sha512_tohex(util__aes_encrypt([]byte(server_config.aes_key), rPass)))[0:50]

		// Inserting the token to database
		database__insert_game_token(nEMID, gamekey)
		util__debug_info(fmt.Sprintf("[WS]| WS-SERVER %s - Token Inserted to DB", client))

		SoForLauncher := gamekey + conf__launcher_gamekey_salt + ";" + file_security_digest//server_config.files_digest
		// Expire the player token
		go func(userID int) {
			<-time.After(59 * time.Second)
			database__delete_game_token(userID)
		}(nEMID)

		// Add IP to firewall
		chan_sentry <- struct_sentry{cmd:"Add", ip:client, rulename: "ArithiaGame"}
		websocket__ws_remote_connectionsRWMutex.Lock()
		websocket__ws_remote_connections = append(websocket__ws_remote_connections, client)
		websocket__ws_remote_connectionsRWMutex.Unlock()

		time.Sleep(2 * time.Second)
		ws.WriteMessage(mt, []byte(SoForLauncher))
		util__debug_info("[MAIN]| WS-SERVER Gamekey sent to " + client)
	}
}

func websocket__delete_conn(sockCli ClientConn, ip_string string){
	websocket__ActiveClientsRWMutex.Lock()
	delete(websocket__ActiveClients, sockCli)
	websocket__ActiveClientsRWMutex.Unlock()
	for i, v := range websocket__ws_remote_connections {
		if(v == ip_string){
			websocket__ws_remote_connectionsRWMutex.Lock()
			websocket__ws_remote_connections = append(websocket__ws_remote_connections[:i], websocket__ws_remote_connections[i+1:]...)
			websocket__ws_remote_connectionsRWMutex.Unlock()
			break
		}
	}
}