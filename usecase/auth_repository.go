package usecase
import "task-manager/domain"


type AuthRepository interface {
    FindUserByUsername(username string) (domain.User, error)
    SaveSession(sessionID string, userID int, expiresAt int64) error
    GetSession(sessionID string) (int, error)
    GetUserRole(userID int) (string, error)
	DeleteSession(sessionID string) error
}
