package appmgr

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// The API uses admin field names; persisted entities retain the app contract.
func bindAppRecord(c *gin.Context, dst any) error {
	var raw map[string]json.RawMessage
	if err := c.ShouldBindJSON(&raw); err != nil {
		return err
	}
	if uid, ok := raw["uid"]; ok {
		raw["uuid"] = uid
	}
	// Identity and timestamps are assigned by the server on creation.
	delete(raw, "createTime")
	delete(raw, "updateTime")
	data, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func presentAppRecord(value any) any {
	data, err := json.Marshal(value)
	if err != nil {
		return value
	}
	var result any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if decoder.Decode(&result) != nil {
		return value
	}
	return renameRecordFields(result)
}

func renameRecordFields(value any) any {
	switch v := value.(type) {
	case []any:
		for i := range v {
			v[i] = renameRecordFields(v[i])
		}
	case map[string]any:
		for _, pair := range [][2]string{{"uuid", "uid"}, {"createTime", "createdAt"}, {"updateTime", "updatedAt"}, {"vipRemainDays", "experienceDeadline"}} {
			if field, ok := v[pair[0]]; ok {
				v[pair[1]] = field
				delete(v, pair[0])
			}
		}
		for k, field := range v {
			v[k] = renameRecordFields(field)
		}
	}
	return value
}
