package models

type AccessDeniedDetail struct {
	RequiredScopes []string `json:"requiredScopes"`
}
