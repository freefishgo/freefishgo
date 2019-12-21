package fishSession

import "time"

type Session struct {
}

func (s Session) GetSession(KeyValue string) (map[string]interface{}, error) {
	return nil, nil
}

func (s Session) GetSessionKeyValue() (string, error) {
	return "", nil
}

func (s Session) SetSession(SessionName string, m map[string]interface{}, duration time.Duration) error {
	return nil
}

func (s Session) UpdateDataTime(SessionName string, duration time.Duration) error {
	return nil
}
