// https://github.com/NotEnoughUpdates/ursa-minor/tree/master

package util

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

var userTokens map[*websocket.Conn]string = make(map[*websocket.Conn]string)

func VerifyUser(username string, serverId string, sessionToken string, conn *websocket.Conn) (bool, string) {
	if userTokens[conn] == sessionToken {
		return true, sessionToken
	}

	var requestUrl = "https://sessionserver.mojang.com/session/minecraft/hasJoined?username=" + username + "&serverId=" + serverId
	res, err := http.Get(requestUrl)
	if err != nil {
		fmt.Print("Error authenticating user: " + err.Error())
		return false, ""
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Print("Error reading body: " + err.Error())
		return false, ""
	}

	var newToken = GenerateNewToken()

	userTokens[conn] = newToken
	return len(body) > 0, newToken
}
