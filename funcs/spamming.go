package funcs

import (
	"context"
	"bytes"
	"net/http"
	"sync"
	"time"
	"io"
	"encoding/json"
	"fmt"
)

func GetChannelIDs(guildID string, token string) ([]string, error) {
	url := fmt.Sprintf("https://discord.com/api/v10/guilds/%s/channels", guildID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bot "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API error: %d %s", res.StatusCode, string(data))
	}

	var channels []struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(res.Body).Decode(&channels); err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(channels))

	for _, channel := range channels {
		ids = append(ids, channel.ID)
	}

	return ids, nil
}

func ChangeGuildName(guildID, newName string, wg *sync.WaitGroup, token string) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://discord.com/api/v10/guilds/%s", guildID)

	payload, err := json.Marshal(map[string]string{
		"name": newName,
	})
	if err != nil {
		fmt.Println(cyan + "-" + reset + " JSON Error")
		return
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		url,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		fmt.Println(cyan + "-" + reset + " Request Error")
		return
	}

	req.Header.Set("Authorization", "Bot "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(cyan + "-" + reset + " HTTP Error")
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		fmt.Println(cyan + "+" + reset + " Guild name changed")
	} else {
		fmt.Printf(cyan+"-"+reset+" Failed: %d\n", res.StatusCode)
	}
}

func DeleteChannel(channelID string, wg *sync.WaitGroup, token string) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s", channelID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		fmt.Println(cyan + "-" + reset + " Request Error")
		return
	}

	req.Header.Set("Authorization", "Bot "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(cyan + "-" + reset + " HTTP Error")
		return
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(cyan + "-" + reset + " Read Error")
		return
	}

	if res.StatusCode == http.StatusNoContent {
		fmt.Println(cyan + "+" + reset + " Channel deleted")
	} else {
		fmt.Printf(cyan+"-"+reset+" Failed: %d %s\n", res.StatusCode, string(data))
	}
}

func PostMessage(guildID, channelName, embed string, wg *sync.WaitGroup, token string) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payload, _ := json.Marshal(map[string]interface{}{
		"name": channelName,
		"type": 0,
	})

	url := fmt.Sprintf("https://discord.com/api/v10/guilds/%s/channels", guildID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bot "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)

	var channel struct {
		ID string `json:"id"`
	}

	if res.StatusCode > 300 || json.Unmarshal(data, &channel) != nil {
		fmt.Println(cyan + "-" + reset + " Channel creation failed")
		return
	}

	msgURL := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channel.ID)

	msgReq, err := http.NewRequestWithContext(ctx, http.MethodPost, msgURL, bytes.NewBuffer([]byte(embed)))
	if err != nil {
		return
	}

	msgReq.Header.Set("Authorization", "Bot "+token)
	msgReq.Header.Set("Content-Type", "application/json")

	msgRes, err := http.DefaultClient.Do(msgReq)
	if err != nil {
		return
	}
	defer msgRes.Body.Close()

	if msgRes.StatusCode == http.StatusOK {
		fmt.Println(cyan + "+" + reset + " Success")
	} else {
		fmt.Println(cyan + "-" + reset + " Message failed")
	}
}
