package api

// TeamMember represents a member of a team
type TeamMember struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	Role      string `json:"role"`
}

// AddMemberRequest represents the request body for adding a team member
type AddMemberRequest struct {
	Member []TeamMember `json:"member"`
	TeamID string       `json:"team_id"`
}

// RemoveMemberRequest represents the request body for removing a team member
type RemoveMemberRequest struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	TeamID    string `json:"team_id"`
}

// TeamInfo represents detailed team information
type TeamInfo struct {
	TeamAlias        string       `json:"team_alias"`
	TeamID           string       `json:"team_id"`
	MembersWithRoles []TeamMember `json:"members_with_roles"`
	Models           []string     `json:"models"`
	Spend            float64      `json:"spend"`
	CreatedAt        string       `json:"created_at"`
}

// TeamResponse represents the API response for team operations
type TeamResponse struct {
	TeamID    string    `json:"team_id"`
	TeamAlias string    `json:"team_alias"`
	TeamInfo  TeamInfo  `json:"team_info"`
	Keys      []KeyInfo `json:"keys"`
}

// KeyInfo represents detailed information about an API key
type KeyInfo struct {
	KeyName   string                 `json:"key_name"`
	KeyAlias  string                 `json:"key_alias"`
	Spend     float64                `json:"spend"`
	Models    []string               `json:"models"`
	TeamID    string                 `json:"team_id"`
	UserID    string                 `json:"user_id"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

// KeyResponse represents the API response for a key info request
type KeyResponse struct {
	Key  string  `json:"key"`
	Info KeyInfo `json:"info"`
}

// KeyListResponse represents the API response for listing keys
type KeyListResponse struct {
	Keys        []string `json:"keys"`
	TotalCount  int      `json:"total_count"`
	CurrentPage int      `json:"current_page"`
	TotalPages  int      `json:"total_pages"`
}

// Team represents detailed team information
type Team struct {
	TeamID    string   `json:"team_id"`
	TeamAlias string   `json:"team_alias"`
	Spend     float64  `json:"spend"`
	Models    []string `json:"models"`
	CreatedAt string   `json:"created_at"`
}

// UserInfo represents detailed user information
type UserInfo struct {
	UserID    string                 `json:"user_id"`
	UserAlias string                 `json:"user_alias"`
	UserEmail string                 `json:"user_email"`
	Teams     []string               `json:"teams"`
	UserRole  string                 `json:"user_role"`
	MaxBudget float64                `json:"max_budget"`
	Spend     float64                `json:"spend"`
	Models    []string               `json:"models"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

// UserResponse represents the API response for user operations
type UserResponse struct {
	UserID   string     `json:"user_id"`
	UserInfo *UserInfo  `json:"user_info"`
	Keys     []KeyInfo  `json:"keys"`
	Teams    []TeamInfo `json:"teams"`
}

// Error represents an API error response
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
