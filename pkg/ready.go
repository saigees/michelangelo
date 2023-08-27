package pkg

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/saigees/michelangelo/internal"
)

func DiscordReady(s *discordgo.Session) {
	s.AddHandler(func( s*discordgo.Session, r *discordgo.Ready) {
		internal.Logger.Infof("Logged in as: %v#%v -> Guild ID: %s", s.State.User.Username, s.State.User.Discriminator, os.Getenv("GUILD_ID"))
	})
}