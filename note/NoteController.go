package note

import (
	"encoding/json"

	"github.com/jh-bate/fantail-bot/store"
)

var noteStore store.Store

func init() {
	noteStore = store.NewRedisStore()
	return
}

func (this *Note) Save() error {
	return noteStore.Save(this.UserId, this)
}

func (this *Note) Delete() error {
	return noteStore.Delete(this.UserId, this)
}

func GetNotes(userid string) (Notes, error) {
	var all Notes
	items, err := noteStore.ReadAll(userid)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(items, &all)

	return all, nil
}
