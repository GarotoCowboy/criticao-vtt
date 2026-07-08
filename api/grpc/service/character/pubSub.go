package character

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/character"
	"github.com/google/uuid"
)

func (c *CharacterService) subscribe(id uint, stream character.CharacterService_UpdateSheetServer) string {
	subID := uuid.NewString()
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.subscribers[id] == nil {
		c.subscribers[id] = make(map[string]character.CharacterService_UpdateSheetServer)
	}
	c.subscribers[id][subID] = stream
	c.Logger.InfoF("new subscriber id:%v subID:%v", id, subID)
	return subID
}

func (c *CharacterService) unsubscribe(id uint, subID string){
	c.mu.Lock()
	defer c.mu.Unlock()

	if subs,ok := c.subscribers[id]; ok{
		delete(subs,subID)
		if len(subs) == 0 {
			delete(c.subscribers, id)
		}
		c.Logger.InfoF("new UnsubscribeToTopic id:%v subID:%v", id, subID)
	} else {
		c.Logger.ErrorF("trying unsubscribe failed with Id:%d, subID:%v", id, subID)
	}
}
