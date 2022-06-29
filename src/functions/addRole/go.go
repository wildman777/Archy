package addRole

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

type Payload struct {
	ServerId string `json:"server_id"`
	UserId   string `json:"user_id"`
	RoleName string `json:"role_name"`
}

// AddRole to as user
func AddRole(w http.ResponseWriter, r *http.Request) {

	// Parse body to get Payload
	var payload = Payload{}
	json.NewDecoder(r.Body).Decode(&payload)

	// Instanciate Discord bot
	dg := instanciateBot()

	// Get the RoleId by matching name
	roleId := getRoleId(dg, payload)

	// Add role to user
	log.Print("Adding role " + payload.RoleName + " to user " + payload.UserId)
	err := dg.GuildMemberRoleAdd(payload.ServerId, payload.UserId, roleId)

	if err != nil {
		panic(err)
	}
	log.Print("Done!")
}

// Instanciate the bot and return the session
func instanciateBot() *discordgo.Session {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	if err != nil {
		error_message := []byte(err.Error())
		error_400_regex, _ := regexp.Compile("400")
		if len(error_400_regex.Find(error_message)) > 0 {
			panic("Can't create Channel - Bad ChannelId")
		}
		error_401_regex, _ := regexp.Compile("401")
		if len(error_401_regex.Find(error_message)) > 0 {
			panic("Unauthorized to create the connection. Verify Discord Token")
		}
		panic(err)
	}
	return dg
}

// Get the roleId by searching the name in existing role
func getRoleId(dg *discordgo.Session, payload Payload) string {

	log.Print(payload.RoleName)

	guildRole, _ := dg.GuildRoles(payload.ServerId)

	for _, v := range guildRole {
		if v.Name == payload.RoleName {
			return v.ID
		}
	}

	panic("Role " + payload.RoleName + " not found")
}
