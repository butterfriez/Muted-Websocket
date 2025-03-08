// https://github.com/NotEnoughUpdates/ursa-minor/tree/master

package util

import (
	"fmt"
	"io"
	"net/http"
)

func VerifyUser(username string, serverId string, sessionToken string) (bool, string) {
	if sessionToken != "" {
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
	}

	var newToken = GenerateNewToken()
	return len(body) > 0, newToken
}
