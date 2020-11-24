package structtomap

import (
	"encoding/json"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
)

func StructToMap(data interface{}) (map[string]interface{}, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	log.GetLogger().WithField("json", string(dataBytes)).Debug("marshalled")

	mapData := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &mapData)
	if err != nil {
		return nil, err
	}

	log.GetLogger().Debugf("unmarshalled: %#v", mapData)

	return mapData, nil
}
