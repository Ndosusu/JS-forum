package directmsg

import (
	"context"
	"errors"
	"fmt"
	"time"

	forumServer "github.com/RealTimeForumDosJean/server"
	localDatabase "github.com/RealTimeForumDosJean/server/database"
)

type DirectMessage struct {
	Msg_id      string    `json:"msg_id"`
	Sender_id   string    `json:"sender_id"`
	Receiver_id string    `json:"receiver_id"`
	Message     string    `json:"message"`
	Sent_at     time.Time `json:"sent_at"`
}

func fromMap(values map[string]any) DirectMessage {
	newDM := DirectMessage{}

	if msgID, ok := values["msg_id"].(string); ok {
		newDM.Msg_id = msgID
	}
	if senderID, ok := values["sender_id"].(string); ok {
		newDM.Sender_id = senderID
	}
	if receiverID, ok := values["receiver_id"].(string); ok {
		newDM.Receiver_id = receiverID
	}
	if msg, ok := values["message"].(string); ok {
		newDM.Message = msg
	}
	if sentAt, ok := values["sent_at"].(time.Time); ok {
		newDM.Sent_at = sentAt
	}

	return newDM
}

func CreateDM(ctx context.Context, params map[string]any) (*DirectMessage, error) {

	msg_ID, _ := forumServer.GenerateUUID()

	sender_ID, senderOK := params["sender_id"].(string)
	receiver_ID, receiverOK := params["receiver_id"].(string)
	message, msgOK := params["message"].(string)

	if !senderOK || !receiverOK || !msgOK {
		return nil, errors.New("informations manquantes pour le message")
	}

	creationDate := time.Now()

	// Insertion du message dans la table
	createMsgQuery := "INSERT INTO directMsgs (msg_id, sender_id, receiver_id, message, sent_at) VALUES (?, ?, ?, ?, ?)"
	_, err := localDatabase.RunDatabaseQuery(ctx, createMsgQuery, msg_ID, sender_ID, receiver_ID, message, creationDate, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du message: %v", err)
	}

	// Création de la structure DirectMessage
	newMsg := &DirectMessage{
		Msg_id:      msg_ID,
		Sender_id:   sender_ID,
		Receiver_id: receiver_ID,
		Message:     message,
		Sent_at:     creationDate,
	}

	return newMsg, nil
}

func FetchMyDMs(ctx context.Context, myUUID string) ([]DirectMessage, error) {
	results, err := localDatabase.RunDatabaseQuery(ctx, "SELECT * FROM directMsgs WHERE sender_uuid = "+myUUID+" OR receiver_uuid = "+myUUID)
	if err != nil {
		return nil, err
	}

	var dmList []DirectMessage

	for _, row := range results {
		dmList = append(dmList, fromMap(row))
	}

	return dmList, nil
}
