package main

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

func getRoleID(guildID string, roleName string) (string, error) {
	guild, err := session.Guild(guildID)
	if err != nil {
		log.Println("getRoleID session.Guild(" + guildID + ") " + err.Error())
		return "", err
	}
	roles := guild.Roles
	for _, role := range roles {
		if role.Name == roleName {
			return role.ID, nil
		}
	}
	return "", errors.New("no " + roleName + " role available")
}

func hasRole(member *discordgo.Member, roleName, guildID string) bool {
	//z jakiegos powodu w strukturze member GuildID jest puste...
	adminRole, err := getRoleID(guildID, roleName)
	if err != nil {
		log.Println("hasRole getRoleID(" + guildID + ", " + roleName + ") " + err.Error())
		return false
	}
	for _, role := range member.Roles {
		if role == adminRole {
			return true
		}
	}
	return false
}
func hasPermission(member *discordgo.Member, guildID string, permission int) bool {
	for _, roleID := range member.Roles {
		role, err := session.State.Role(guildID, roleID)
		if err != nil {
			log.Println("hasPermisson session.State.Role(" + guildID + ", " + roleID + ") " + err.Error())
			return false
		}
		if role.Permissions&permission != 0 {
			return true
		}
	}
	return false
}

func hasAdminPermissions(member *discordgo.Member, guildID string) bool {
	if hasRole(member, getAdminRoleForGuild(guildID), guildID) || hasPermission(member, guildID, 8) { // 8 - administrator
		return true
	}
	return false
}

func getAdminRoleForGuild(guildID string) string {
	var serverConfig ServerConfig
	err := DbMap.SelectOne(&serverConfig, "SELECT * FROM ServerConfig")
	if err != nil {
		log.Println("getAdminRoleForGuild(" + guildID + ") " + err.Error())
		return ""
	}
	return serverConfig.AdminRole
}