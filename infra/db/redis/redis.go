package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"boilerplate-go/data/protocols"
	"github.com/go-redis/redis"
)

type RedisHelper struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) protocols.Redis {
	return &RedisHelper{client}
}

func (helper RedisHelper) FindAll() ([]protocols.FlagPayload, error) {
	pattern := "rollout:*"
	// Recuperar a estrutura existente do Redis
	keys, err := helper.client.Keys(pattern).Result()
	if err != nil {
		return []protocols.FlagPayload{}, err
	}

	var payloads []protocols.FlagPayload

	for _, key := range keys {
		// Remover o prefixo "rollout:" de cada chave
		keyWithoutPrefix := strings.TrimPrefix(key, "rollout:")

		// Recuperar o valor da chave
		value, err := helper.client.Get(key).Result()
		if err != nil {
			return []protocols.FlagPayload{}, err
		}

		// Deserializar o valor para FlagPayload
		var payload protocols.FlagPayload
		err = json.Unmarshal([]byte(value), &payload)
		if err != nil {
			return []protocols.FlagPayload{}, err
		}

		// Criar array com a quantidade de IDs
		idsLen := []string{strconv.Itoa(len(payload.IDs))}
		// Remover os IDs
		cleanPayload := protocols.FlagPayload{
			Flagger:      keyWithoutPrefix,
			ExpirationAt: payload.ExpirationAt,
			FullRollout:  payload.FullRollout,
			Description:  payload.Description,
			IDs:          idsLen,
		}

		payloads = append(payloads, cleanPayload)
	}

	return payloads, nil
}

func (helper RedisHelper) FindFlagger(flagger string) (protocols.FlagPayload, error) {
	// Recuperar a estrutura existente do Redis
	existingJSON, err := helper.client.Get("rollout:" + flagger).Result()
	if err != nil {
		if err == redis.Nil {
			return protocols.FlagPayload{}, fmt.Errorf("flagger '%s' is not exist", flagger)
		}
		return protocols.FlagPayload{}, err
	}

	// Converter o JSON recuperado em uma estrutura FlaggerPayload
	var existingPayload protocols.FlagPayload
	if existingJSON != "" {
		err = json.Unmarshal([]byte(existingJSON), &existingPayload)
		if err != nil {
			return protocols.FlagPayload{}, err
		}
	}
	existingPayload.Flagger = strings.TrimPrefix(existingPayload.Flagger, "rollout:")
	return existingPayload, nil
}

func (helper RedisHelper) FindFlaggersCompany(id string) ([]string, error) {
	pattern := "rollout:*"
	// Busca todas as flaggers cadastrada no Redis
	keys, err := helper.client.Keys(pattern).Result()
	if err != nil {
		return nil, err
	}

	var flaggers []string
	for _, key := range keys {
		// Pega o valor associado à chave do rollout
		valueFlagger, err := helper.client.Get(key).Result()
		if err != nil {
			return nil, err
		}

		// Converter o JSON em uma estrutura FlaggerPayload
		var payload protocols.FlagPayload
		if valueFlagger != "" {
			err = json.Unmarshal([]byte(valueFlagger), &payload)
			if err != nil {
				return nil, err
			}
		}

		// Verificar se a company está participando desse rollout ou se estar full rollout
		if payload.FullRollout || contains(payload.IDs, id) {
			// Remover o prefixo "rollout:" da chave
			flaggers = append(flaggers, strings.TrimPrefix(key, "rollout:"))
		}
	}

	if flaggers != nil {
		// Faz o retorno do array contendo as flaggers de rollouts
		return flaggers, nil
	} else {
		// Caso não tenha nenhum rollout é rotornado um array vazio
		return []string{}, nil
	}
}

func (helper RedisHelper) CreateFlagger(payload protocols.FlagPayload) error {
	// Verificar se a chave já existe
	exists, err := helper.client.Exists("rollout:" + payload.Flagger).Result()
	if err != nil {
		return err
	}

	// Se a chave não existe, então cadastramos
	if exists == 0 {
		payload.Flagger = "rollout:" + payload.Flagger
		// Convertendo o payload para JSON
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		// Definindo a chave com o payload como valor e sem data de expiração
		err = helper.client.Set(payload.Flagger, payloadJSON, 0).Err()
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("flagger '%s' already exist", payload.Flagger)
	}
}

func (helper RedisHelper) DeleteFlagger(flagger string) error {
	deleted, err := helper.client.Del("rollout:" + flagger).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return fmt.Errorf("key not found for flagger: %s", flagger)
	}
	return nil
}

func (helper RedisHelper) AddCompanyInFlagger(flagger string, ids []string) error {
	// Recuperar a estrutura existente do Redis
	existingJSON, err := helper.client.Get("rollout:" + flagger).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("flagger '%s' is not exist", flagger)
		}
		return err
	}

	// Converter o JSON recuperado em uma estrutura FlaggerPayload
	var existingPayload protocols.FlagPayload
	if existingJSON != "" {
		err = json.Unmarshal([]byte(existingJSON), &existingPayload)
		if err != nil {
			return err
		}
	}

	// Mapear os IDs existentes para uma estrutura de dados para rápida verificação de existência
	existingIDsMap := make(map[string]struct{})
	for _, id := range existingPayload.IDs {
		existingIDsMap[id] = struct{}{}
	}

	// Adicionar apenas os IDs que ainda não estão presentes na lista existente
	var newIDs []string
	for _, id := range ids {
		if _, exists := existingIDsMap[id]; !exists {
			newIDs = append(newIDs, id)
		}
	}

	// Se não houver novos IDs a serem adicionados, retornar imediatamente
	if len(newIDs) == 0 {
		return nil
	}

	// Adicionar os novos IDs à lista existente
	existingPayload.IDs = append(existingPayload.IDs, newIDs...)

	// Converter a estrutura atualizada de volta para JSON
	updatedJSON, err := json.Marshal(existingPayload)
	if err != nil {
		return err
	}

	// Atualizar o valor no Redis com a estrutura atualizada
	err = helper.client.Set("rollout:"+flagger, updatedJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (helper RedisHelper) RemoveCompanyInFlagger(flagger string, ids []string) error {
	// Recuperar a estrutura existente do Redis
	existingJSON, err := helper.client.Get("rollout:" + flagger).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("flagger '%s' is not exist", flagger)
		}
		return err
	}

	// Converter o JSON recuperado em uma estrutura FlaggerPayload
	var existingPayload protocols.FlagPayload
	if existingJSON != "" {
		err = json.Unmarshal([]byte(existingJSON), &existingPayload)
		if err != nil {
			return err
		}
	}

	// Remover os IDs da lista na estrutura
	idsToRemoveMap := make(map[string]struct{})
	for _, id := range ids {
		idsToRemoveMap[id] = struct{}{}
	}

	var updatedIDs []string
	for _, id := range existingPayload.IDs {
		if _, exists := idsToRemoveMap[id]; !exists {
			updatedIDs = append(updatedIDs, id)
		}
	}
	existingPayload.IDs = updatedIDs

	// Converter a estrutura atualizada de volta para JSON
	updatedJSON, err := json.Marshal(existingPayload)
	if err != nil {
		return err
	}

	// Atualizar o valor no Redis com a estrutura atualizada
	err = helper.client.Set("rollout:"+flagger, updatedJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (helper RedisHelper) UpdateFullRollout(flagger string, status bool) error {
	// Recuperar a estrutura existente do Redis
	existingJSON, err := helper.client.Get("rollout:" + flagger).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("flagger '%s' is not exist", flagger)
		}
		return err
	}

	// Converter o JSON recuperado em uma estrutura FlaggerPayload
	var existingPayload protocols.FlagPayload
	if existingJSON != "" {
		err = json.Unmarshal([]byte(existingJSON), &existingPayload)
		if err != nil {
			return err
		}
	}

	existingPayload.FullRollout = status

	// Converter a estrutura atualizada de volta para JSON
	updatedJSON, err := json.Marshal(existingPayload)
	if err != nil {
		return err
	}

	// Atualizar o valor no Redis com a estrutura atualizada
	err = helper.client.Set("rollout:"+flagger, updatedJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// Função auxiliar para verificar se um elemento está presente em uma slice de strings
func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
