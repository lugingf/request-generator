package resources

//import (
//	"context"
//	_ "github.com/jinzhu/gorm/dialects/mysql"
//	"stash.tutu.ru/golang/log"
//	"stash.tutu.ru/golang/resources/mongo"
//)
//
//func (r *Resources) initMongo(ctx context.Context) error {
//
//	collection, err := mongo.Collection(ctx, mongo.Config{})
//	if err != nil {
//		return err
//	}
//
//	r.Mongo = collection
//
//	log.Logger.Info().Msg("initMongo success")
//	return nil
//}
