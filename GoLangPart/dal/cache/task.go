package cache

import (
	"log"
)

func CheckTaskExists(itemId int64) (bool, error) {
	conn, err := GetRedisConn()
	if err != nil {
		log.Printf("get redis conn err: %v",err)
		return false, err
	}
	defer conn.Close()

	resp, err := conn.Do("SISMEMBER", "TaskSet", itemId)
	if err != nil {
		log.Printf("redis call SISMEMBER TaskSet %v err: %v", itemId, err)
		return false, err
	}
	return resp.(int64) == 1, err
}

func AddTaskToCache(itemId int64) error {
	conn, err := GetRedisConn()
	if err != nil {
		log.Printf("get redis conn err: %v",err)
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SADD", "TaskSet", itemId)
	if err != nil {
		log.Printf("redis call SADD TaskSet %v err: %v", itemId, err)
		return err
	}
	return nil
}

func DeleteTaskById(taskId int64) error {
	conn, err := GetRedisConn()
	if err != nil {
		log.Fatalf("get redis conn err: %v", err)
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SREM", "TaskSet", taskId)
	if err != nil {
		log.Printf("redis call SREM TaskSet %v err: %v", taskId, err)
		return err
	}
	return nil

}