package auth

type AuthenticatedUser struct {
	CompanyID string `json:"companyId"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	UserID    string `json:"userId"`
}
