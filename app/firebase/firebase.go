package firebase

import (
	"context"
	"errors"
	"net/http"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var (
	firebaseApp *firebase.App
)

// InitFirebase initializes the Firebase Admin SDK with the provided service account JSON file path.
func InitFirebase(serviceAccountFilePath string) error {
	opt := option.WithCredentialsFile(serviceAccountFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}
	firebaseApp = app
	return nil
}

// VerifyTokenAndGetUserID verifies the Firebase JWT token from the request and returns the user's ID.
func VerifyTokenAndGetUserID(r *http.Request) (string, error) {
	if firebaseApp == nil {
		return "", errors.New("Firebase app is not initialized")
	}

	token := extractTokenFromHeader(r)
	if token == "" {
		return "", errors.New("Authorization token not found")
	}

	ctx := context.Background()
	client, err := firebaseApp.Auth(ctx)
	if err != nil {
		return "", err
	}

	decodedToken, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		return "", err
	}

	return decodedToken.UID, nil
}

// extractTokenFromHeader extracts the JWT token from the Authorization header.
func extractTokenFromHeader(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header != "" {
		// Assuming the token is in the format "Bearer <token>"
		return header[len("Bearer "):]
	}
	return ""
}
