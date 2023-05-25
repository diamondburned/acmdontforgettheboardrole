package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("Missing BOT_TOKEN environment variable")
	}

	s := state.New(token)
	s.AddIntents(gateway.IntentGuildMembers)
	s.AddHandler(func(ev *gateway.GuildMemberUpdateEvent) {
		onRoleChange(s, Member{
			GuildID: ev.GuildID,
			User:    ev.User,
			Nick:    ev.Nick,
			RoleIDs: ev.RoleIDs,
			Avatar:  ev.Avatar,
		})
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := s.Connect(ctx); err != nil {
		log.Fatalln("Failed to connect:", err)
	}
}

func onRoleChange(s *state.State, member Member) {
	check, err := checkMemberRole(s, member)
	if err != nil {
		log.Println("Failed to check member role:", err)
		return
	}
	updateMemberRole(s, member, check)
}
