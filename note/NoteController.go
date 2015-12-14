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

//todo this need to be moved out
func bytes(v interface{}) []byte {
	switch v := v.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	}
	return nil
}

func GetNotes(userid string) (Notes, error) {
	var all Notes
	items, err := noteStore.ReadAll(userid)
	if err != nil {
		return nil, err
	}

	for i := range items {
		var n Note
		json.Unmarshal(bytes(items[i]), &n)
		all = append(all, &n)
	}

	return all, nil
}
