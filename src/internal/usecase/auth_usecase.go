package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	//"log"
	"task-manager/src/internal/adaptors/persistance"
	//"task-manager/src/internal/core/user"

	"time"

	"golang.org/x/crypto/bcrypt"
)


type AuthUseCase struct {
   userRepo *persistance.UserPostgresRepo
   sessionRepo *persistance.SessionPostgresRepo
}

func NewAuthUseCase(userRepo *persistance.UserPostgresRepo, sessionRepo *persistance.SessionPostgresRepo) *AuthUseCase {
   return &AuthUseCase{
    userRepo: userRepo,
    sessionRepo: sessionRepo,
   }
}

func (uc *AuthUseCase) Login(username, password string) (string, error) {
    user, err := uc.userRepo.FindUserByUsername(username)
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
    err = uc.sessionRepo.SaveSession(sessionID, user.ID, expiresAt)
    if err != nil {
        return "", err
    }
   // log.Println("Alright...!") --- Debug
    return sessionID, nil
}

func (uc *AuthUseCase) Authorize(sessionID, requiredRole string) (bool, error) {
    userID, err := uc.sessionRepo.GetSession(sessionID)
    if err != nil {
        return false, err
    }
    role, err := uc.userRepo.GetUserRole(userID)
    if err != nil {
        return false, err
    }
    return role == requiredRole || role == "admin", nil // Admin overrides all
}

func (uc *AuthUseCase) Logout(sessionID string) error{
	return uc.sessionRepo.DeleteSession(sessionID)
}

func generateSessionID() (string, error) {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
}