package firebase

import (
	"context"
	"encoding/base64"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
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

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}

	//Messaging client
	s, err := app.Messaging(ctx)

	f = &Firebase{
		Client: s,
	}

	return
}

func (f *Firebase) Broadcast(ctx context.Context, m *messaging.MulticastMessage, dc int64) (success bool, err error) {
	r, err := f.Client.SendMulticast(ctx, m)

	success = int64(r.SuccessCount) == dc

	return
}
