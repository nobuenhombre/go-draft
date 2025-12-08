package vars

import (
	"errors"
	"strings"

	"github.com/nobuenhombre/suikat/pkg/ge"
)

var (
	ErrorInvalidKeyValuePair = errors.New("invalid key/value pair")
	ErrorEmptyKey            = errors.New("empty key")
)

type Service interface {
	Parse(vars string) (map[string]string, error)
}

type Provider struct {
}

func New() Service {
	return &Provider{}
}

func (p *Provider) Parse(vars string) (map[string]string, error) {
	result := make(map[string]string)

	if strings.TrimSpace(vars) == "" {
		return result, nil
	}

	// Разделяем строку по запятым на пары ключ-значение
	pairs := strings.Split(vars, ",")

	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)

		if len(kv) != 2 {
			return nil, ge.Pin(ErrorInvalidKeyValuePair, ge.Params{"pair": pair})
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		if key == "" {
			return nil, ge.Pin(ErrorEmptyKey, ge.Params{"pair": pair})
		}

		result[key] = value
	}

	return result, nil
}
