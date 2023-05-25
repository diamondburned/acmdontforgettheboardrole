package main

import (
	"fmt"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

const acmGuildID = 710225099923521558

var boardRoles = setFromList([]string{
	"AI Board",
	"Algo Board",
	"Design Board",
	"Dev Board",
	"Game Dev Board",
	"Special Events Board",
	"Marketing Board",
	"Node Buds Board",
})

var (
	boardRoleName = "Board"
	boardRoleID   = discord.RoleID(710225414706036746) // act as guard
)

type Member struct {
	GuildID discord.GuildID
	User    discord.User
	Nick    string
	RoleIDs []discord.RoleID
	Avatar  discord.Hash
}

func checkMemberRole(s *state.State, member Member) (bool, error) {
	if acmGuildID != member.GuildID {
		return false, nil
	}

	roles, err := s.Roles(member.GuildID)
	if err != nil {
		return false, fmt.Errorf("failed to get roles: %w", err)
	}

	boardRole := getBoardRole(roles)
	if boardRole == nil || boardRole.Name != boardRoleName {
		return false, fmt.Errorf("missing board role, please reconfigure bot")
	}

	memberRoles := filterMemberRoles(roles, setFromList(member.RoleIDs))
	for _, role := range memberRoles {
		// Board-specific roles are always higher than the board role.
		// Ignore the ones that aren't.
		if role.Position < boardRole.Position {
			log.Printf(
				"ignoring role %s (%v) with position %d which is below Board",
				role.Name, role.ID, role.Position,
			)
			continue
		}
		if boardRoles[role.Name] {
			return true, nil
		}
	}

	return false, nil
}

func updateMemberRole(s *state.State, member Member, passes bool) error {
	memberRoles := setFromList(member.RoleIDs)
	if passes {
		if !memberRoles[boardRoleID] {
			return s.AddRole(member.GuildID, member.User.ID, boardRoleID, api.AddRoleData{
				AuditLogReason: "member has board role",
			})
		}
		return nil
	} else {
		if memberRoles[boardRoleID] {
			return s.RemoveRole(member.GuildID, member.User.ID, boardRoleID,
				"member does not have board role",
			)
		}
		return nil
	}
}

func getBoardRole(roles []discord.Role) *discord.Role {
	for i, role := range roles {
		if role.ID == boardRoleID {
			return &roles[i]
		}
	}
	return nil
}

func filterMemberRoles(roles []discord.Role, memberRoleIDs Set[discord.RoleID]) []discord.Role {
	filtered := make([]discord.Role, 0, len(memberRoleIDs))
	for _, role := range roles {
		if _, ok := memberRoleIDs[role.ID]; ok {
			filtered = append(filtered, role)
		}
	}
	return filtered
}

type Set[T comparable] map[T]bool

func setFromList[T comparable](slice []T) Set[T] {
	set := make(map[T]bool, len(slice))
	for _, item := range slice {
		set[item] = true
	}
	return set
}
