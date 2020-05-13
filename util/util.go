package util

import (
	"encoding/json"

	"github.com/gamezop/model"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/gomodule/redigo/redis"
	nsq "github.com/nsqio/go-nsq"
)

//PushDataToRedis : Push data to Redis
func PushDataToRedis(req *gin.Context) (err error) {
	var (
		reqPayload model.RequestModel
		byteArr    []byte
		conn       redis.Conn
	)
	if err = req.ShouldBindJSON(&reqPayload); err == nil {
		if conn, err = redis.Dial("tcp", model.RedisConnect); err == nil {
			defer conn.Close()
			if _, err = conn.Do("HMSET", "data:"+reqPayload.EmpID,
				"name", reqPayload.Name,
				"age", reqPayload.Age,
				"hobby", reqPayload.Hobby); err == nil {
				if byteArr, err = json.Marshal(reqPayload); err == nil {
					//Publish Event To NSQ for syncing data to database
					err = publishToNSQ(byteArr)
				}
			}
		}
	}
	return
}

//GetDataFromDatabase : Get Data From Database
func GetDataFromDatabase(req *gin.Context) (data interface{}, err error) {
	var (
		responseModel model.ResponseModel
		conn          *pg.DB
	)
	if conn, err = pgConn(); err == nil {
		selectSQL := `
			SELECT emp_id, name, age, hobby 
			FROM gamezop;`
		responseModel.Data = make([]model.RequestModel, 0)
		if _, err = conn.Query(&responseModel.Data, selectSQL); err == nil {
			data = responseModel
		}
	}
	return
}

//pgConn : Return PG connection
func pgConn() (conn *pg.DB, err error) {
	conn = pg.Connect(&pg.Options{
		Addr:     model.DBHost + ":" + model.DBPort,
		User:     model.DBUser,
		Database: model.DBDatabase,
		Password: model.DBPassword,
	})
	return
}

//PushEventToDatabase : Push Event to Database
func PushEventToDatabase(body []byte) (err error) {
	var (
		reqModel model.RequestModel
		conn     *pg.DB
	)
	if err = json.Unmarshal(body, &reqModel); err == nil {
		if conn, err = pgConn(); err == nil {
			insertSQL := `
				INSERT INTO gamezop(emp_id, name, age, hobby)
				VALUES(?,?,?,?);`
			params := []interface{}{reqModel.EmpID, reqModel.Name, reqModel.Age, reqModel.Hobby}
			if _, err = conn.Exec(insertSQL, params...); err == nil {
				//Delete data from Redis after inserting into database
				err = deleteFromRedis(reqModel.EmpID)
			}
		}
	}
	return
}

//deleteFromRedis : Delete From Redis
func deleteFromRedis(empID string) (err error) {
	var conn redis.Conn
	if conn, err = redis.Dial("tcp", model.RedisConnect); err == nil {
		_, err = conn.Do("DEL", "data:"+empID)
	}
	return
}

//publishToNSQ : Publish event to NSQ
func publishToNSQ(data []byte) (err error) {
	config := nsq.NewConfig()
	var w *nsq.Producer
	if w, err = nsq.NewProducer("127.0.0.1:4150", config); err == nil {
		err = w.Publish(model.NsqTopic, data)
		w.Stop()
	}
	return
}
