package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	//"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)


type AuthUseCase struct {
    repo AuthRepository
}

func NewAuthUseCase(repo AuthRepository) *AuthUseCase {
    return &AuthUseCase{repo: repo}
}

func (uc *AuthUseCase) Login(username, password string) (string, error) {
    user, err := uc.repo.FindUserByUsername(username)
    if err != nil {
       
	        return "", errors.New("invalid credentials")
    }
 
	if err:= bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "",errors.New("invalid credentials")
	}

    sessionID, err := generateSessionID()
    if err != nil {
        return "", err
    }
    expiresAt := time.Now().Add(1 * time.Hour).Unix()
    err = uc.repo.SaveSession(sessionID, user.ID, expiresAt)
    if err != nil {
        return "", err
    }
   // log.Println("Alright...!") --- Debug
    return sessionID, nil
}

func (uc *AuthUseCase) Authorize(sessionID, requiredRole string) (bool, error) {
    userID, err := uc.repo.GetSession(sessionID)
    if err != nil {
        return false, err
    }
    role, err := uc.repo.GetUserRole(userID)
    if err != nil {
        return false, err
    }
    return role == requiredRole || role == "admin", nil // Admin overrides all
}

func (uc *AuthUseCase) Logout(sessionID string) error{
	return uc.repo.DeleteSession(sessionID)
}

func generateSessionID() (string, error) {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
}