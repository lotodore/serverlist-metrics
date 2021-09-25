package webrequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const (
	// Use GetServerList Steam API, which is not officially documented.
	// See e.g. https://gist.github.com/Decicus/5d6eb057da4b5f228a4f1c2334ce1e2f
	steamGetServerListRequest = "https://api.steampowered.com/IGameServersService/GetServerList/v1/?key=%s&format=json&filter=%s&limit=%d"
	defaultRequestLimit       = 5000
)

// Single game server from steam server response.
type GameServer struct {
	// Names equal json names.
	Addr       string
	Gameport   int
	Steamid    string
	Name       string
	Appid      int
	Gamedir    string
	Version    string
	Product    string
	Region     int
	Players    int
	MaxPlayers int
	Bots       int
	Map        string
	Secure     bool
	Dedicated  bool
	Os         string
	Gametype   string
}

// Server list from steam server response.
type GameServerList struct {
	// Names equal json names.
	Response struct {
		Servers []GameServer
	}
}

// Request game servers using steam web api.
func RequestGameServers(filters ...Filter) ([]GameServer, error) {
	filterString := CreateFilterString(filters)
	if filterString == "" {
		return nil, errors.New("empty server filter should not be used")
	}

	apiKey := os.Getenv("STEAM_WEB_API_KEY")
	if apiKey == "" {
		return nil, errors.New("steam api key not set")
	}

	response, err := http.Get(fmt.Sprintf(
		steamGetServerListRequest, url.QueryEscape(apiKey), url.QueryEscape(filterString), defaultRequestLimit))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var serverList GameServerList
	if err := json.NewDecoder(response.Body).Decode(&serverList); err != nil {
		return nil, err
	}
	return serverList.Response.Servers, nil
}
