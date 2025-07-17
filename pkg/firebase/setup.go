package firebase

import (
	"context"
	"encoding/base64"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Firebase struct {
	Client *messaging.Client
}

func getDecodedFireBaseKey(a string) (d []byte, err error) {
	d, err = base64.StdEncoding.DecodeString(a)
	if err != nil {
		return
	}
	return
}

func Connection(authkey string) (f *Firebase, err error) {
	ctx := context.Background()

	decodedKey, err := getDecodedFireBaseKey(authkey)
	if err != nil {
		return
	}

	opt := option.WithCredentialsJSON(decodedKey)

	// Firebase admin SDK initialization
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err // Return error instead of panic
	}

	// Messaging client
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	f = &Firebase{
		Client: client,
	}

	return
}

func (f *Firebase) Broadcast(ctx context.Context, m *messaging.MulticastMessage, dc int64) (success bool, err error) {
	r, err := f.Client.SendEachForMulticast(ctx, m)
	if err != nil {
		return false, err
	}

	success = int64(r.SuccessCount) == dc
	return
}
