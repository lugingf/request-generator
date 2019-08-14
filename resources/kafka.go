package resources

//import (
//	"context"
//	"stash.tutu.ru/golang/log"
//	"stash.tutu.ru/golang/resources/kafka"
//)
//
//func (r *Resources) initKafka(ctx context.Context) error {
//
//	// consumer
//
//	consumer := kafka.NewResource(ctx)
//	if err := consumer.Init(); err != nil {
//		return err
//	}
//
//	r.KafkaConsumer = consumer
//	log.Logger.Info().Msg("initKafkaConsumer success")
//
//	// producer
//
//	produce := kafka.NewResource(ctx)
//	if err := produce.Init(); err != nil {
//		return err
//	}
//
//	producer, err := produce.GetSyncProducer()
//	if err != nil {
//		return err
//	}
//
//	r.KafkaProducer = producer
//	log.Logger.Info().Msg("initKafkaProducer success")
//
//	return nil
//}
