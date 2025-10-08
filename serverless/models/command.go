package models

type Command struct {
	Name        string          `json:"name" validate:"required"`
	Type        int             `json:"type" validate:"required"`
	Description string          `json:"description" validate:"required"`
	Options     []CommandOption `json:"options,omitempty"`
}

type CommandOption struct {
	Name        string                `json:"name" validate:"required"`
	Description string                `json:"description" validate:"required"`
	Type        int                   `json:"type" validate:"required"`
	Required    bool                  `json:"required,omitempty"`
	Choices     []CommandOptionChoice `json:"choices,omitempty"`
}

type CommandOptionChoice struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}
