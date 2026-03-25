package dto

type LoginRequestDTO struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AuthUserResponseDTO struct {
	User UserSummaryDTO `json:"user"`
}

type LogoutResponseDTO struct {
	OK bool `json:"ok"`
}
