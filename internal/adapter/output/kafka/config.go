package kafka

// Config contém as configurações para conexão com Kafka.
type Config struct {
	Brokers []string // Lista de brokers (ex: ["localhost:9092"])
	Topic   string   // Nome do tópico
	GroupID string   // ID do grupo de consumidores (apenas para consumer)
}
