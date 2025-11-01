package customer

type RegisterDto struct {
	Name     string
	Password string
	Phone    string
}

type LoginDto struct {
	Phone    string
	Password string
}

type VerifyOtpDto struct {
	ChallengeID string
	Code        string
}
