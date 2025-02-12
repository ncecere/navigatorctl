package team

import (
	"github.com/spf13/viper"
)

// Resolver handles the resolution of team identifiers (IDs or aliases)
type Resolver struct {
	aliases map[string]string
}

// NewResolver creates a new team resolver instance
func NewResolver() *Resolver {
	return &Resolver{
		aliases: make(map[string]string),
	}
}

// LoadAliases loads team aliases from the configuration
func (r *Resolver) LoadAliases() error {
	aliases := viper.GetStringMapString("team.aliases")
	if aliases != nil {
		r.aliases = aliases
	}
	return nil
}

// ResolveTeam takes either a team ID or alias and returns the actual team ID
func (r *Resolver) ResolveTeam(identifier string) (string, error) {
	// First check if it's an alias
	if teamID, exists := r.aliases[identifier]; exists {
		return teamID, nil
	}

	// If it's not an alias, assume it's a team ID
	// TODO: Validate team ID format when we have API specs
	return identifier, nil
}

// GetAlias returns the alias for a team ID if one exists
func (r *Resolver) GetAlias(teamID string) string {
	for alias, id := range r.aliases {
		if id == teamID {
			return alias
		}
	}
	return ""
}

// AddAlias adds or updates a team alias
func (r *Resolver) AddAlias(alias, teamID string) {
	r.aliases[alias] = teamID
}

// RemoveAlias removes a team alias
func (r *Resolver) RemoveAlias(alias string) {
	delete(r.aliases, alias)
}
