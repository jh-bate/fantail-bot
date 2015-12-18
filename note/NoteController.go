package note

import (
	"encoding/json"
	"fmt"

	"github.com/jh-bate/fantail-bot/store"
)

var noteStore store.Store

func init() {
	noteStore = store.NewRedisStore()
	return
}

func (this *Note) Save() error {
	return noteStore.Save(fmt.Sprintf("%d", this.UserId), this)
}

func (this *Note) Delete() error {
	return noteStore.Delete(fmt.Sprintf("%d", this.UserId), this)
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
