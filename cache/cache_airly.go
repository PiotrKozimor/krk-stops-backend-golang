package cache

import (
	"context"
	"log"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

var AirlyExpire = time.Minute * 10
var ctx = context.Background()

const AirlyPrefix = "airly-"

func IsAirlyCached(c *redis.Client, installation *pb.Installation) (cached bool, err error) {
	var exist int64
	exist, err = c.Exists(ctx, AirlyPrefix+string(installation.Id)).Result()
	if err != nil {
		return
	}
	if err != nil || exist == 0 {
		cached = false
	} else {
		cached = true
	}
	return
}

func getAirlyKey(i *pb.Installation) string {
	return AirlyPrefix + string(i.Id)
}
func CacheAirly(c *redis.Client, airly *pb.Airly, installation *pb.Installation) (err error) {
	pipe := c.Pipeline()
	executePipe := true
	pipe.Del(ctx, getAirlyKey(installation))
	rawAirly, err := proto.Marshal(airly)
	if err != nil {
		return err
	}
	pipe.Set(ctx, getAirlyKey(installation), string(rawAirly), 0)
	pipe.Expire(ctx, getAirlyKey(installation), AirlyExpire)
	if executePipe {
		_, err = pipe.Exec(ctx)
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}

func GetCachedAirly(c *redis.Client, installation *pb.Installation) (*pb.Airly, error) {
	var airly pb.Airly
	airlyRaw, err := c.Get(ctx, getAirlyKey(installation)).Result()
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal([]byte(airlyRaw), &airly)
	return &airly, err
}
