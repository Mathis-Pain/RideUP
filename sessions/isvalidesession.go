package sessions

import (
	"fmt"
	"time"
)

func IsValidSession(sessionID string) bool {
	session, err := GetSession(sessionID)
	if err != nil {
		fmt.Println("Session non trouvée")
		return false
	}

	// Vérifier que UserID est valide
	if session.UserID == 0 {
		fmt.Println("UserID invalide")
		return false
	}

	// Vérifier que la session n’est pas expirée
	if session.ExpiresAt.Before(time.Now()) {
		fmt.Println("Session expirée")
		return false
	}

	return true
}
