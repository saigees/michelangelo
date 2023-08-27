package michelangelo

import (
	"fmt"
	"os"
	"os/signal"

	"flag"

	"github.com/bwmarrin/discordgo"
	"github.com/saigees/michelangelo/cmd"
	"github.com/saigees/michelangelo/internal"
	"github.com/saigees/michelangelo/pkg"
)

var (
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutting down or not.")
)

var discord *discordgo.Session

func Main() {
	flag.Parse()

	var err error
	discord, err = discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_TOKEN")))
	if err != nil {
		internal.Logger.Error("Failed to create discord client.")
	}

	commands := []func() (*discordgo.ApplicationCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate)){
		cmd.PingCmd,
	}

	cmdDatas := []*discordgo.ApplicationCommand{}
	cmdMap := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
	for _, fn := range commands {
		cmd, fun := fn()
		cmdMap[cmd.Name] = fun
		cmdDatas = append(cmdDatas, cmd)
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := cmdMap[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	pkg.DiscordReady(discord)
	err = discord.Open()
	if err != nil {
		internal.Logger.Error("Failed to open discord session.")
	}
	internal.Logger.Info("Adding commands....")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range cmdDatas {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, os.Getenv("GUILD_ID"), v)
		if err != nil {
			internal.Logger.Errorf("Failed to create command -> %v", v.Name)
		}
		registeredCommands[i] = cmd
	}

	defer discord.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	internal.Logger.Error("Press CTRL+C to exit.")
	<-stop

	if *RemoveCommands {
		internal.Logger.Error("Removing commands....")
		for _, v := range registeredCommands {
			err := discord.ApplicationCommandDelete(discord.State.User.ID, os.Getenv("GUILD_ID"), v.ID)
			if err != nil {
				internal.Logger.Errorf("Failed to delete command -> %v", v.Name)
			}
		}
	}

	internal.Logger.Error("Gracefully shutting down..")
}
