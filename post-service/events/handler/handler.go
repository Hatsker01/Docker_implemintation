package handler

import (
	"github.com/Hatsker01/Docker_implemintation/post-service/config"
	pb "github.com/Hatsker01/Docker_implemintation/post-service/genproto"
	"github.com/Hatsker01/Docker_implemintation/post-service/pkg/logger"
	"github.com/Hatsker01/Docker_implemintation/post-service/storage"
)

type EventHandler struct {
	config  config.Config
	storage storage.IStorage
	log     logger.Logger
}

func NewEventHandlerFunc(config config.Config, storage storage.IStorage, log logger.Logger) *EventHandler {
	return &EventHandler{
		config:  config,
		storage: storage,
		log:     log,
	}
}

func (h *EventHandler) Handler(value []byte) error {
	var user pb.User
	err := user.Unmarshal(value)
	if err != nil {
		return err
	}
	for _, post := range user.Posts {
		_, err := h.storage.Post().CreatePost(post)
		if err != nil {
			return err
		}
	}
	return nil
}
