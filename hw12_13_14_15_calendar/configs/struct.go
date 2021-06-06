package configs

import (
	"time"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage"
)

// Config При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	DB     DBConf
	Server Server
}

type LoggerConf struct {
	Level string
	Path  string
}

type Server struct {
	HTTP string
	Grpc string
}

type DBConf struct {
	User     string
	Password string
	Host     string
	Port     uint64
	Name     string
	Mem      bool
}

type MQ struct {
	URI          string
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	Body         string
	Reliable     bool
	Queue        string
	Tag          string
	Interval     time.Duration
}

type MQNotification struct {
	EventID string
	Title   string
	Date    time.Time
	UserID  string
}

type Sheduler struct {
	Logger LoggerConf
	DB     DBConf
	MQ     MQ
}

type Sender struct {
	Logger LoggerConf
	DB     DBConf
	MQ     MQ
}

func ConvertToMQNotification(event storage.Event) MQNotification {
	return MQNotification{
		EventID: event.ID,
		Title:   event.Title,
		Date:    event.StartData,
		UserID:  event.OwnerID,
	}
}
