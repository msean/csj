package main

import (
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"
)

func main() {

	numWokers := 100
	antsPool, _ := ants.NewPool(numWokers)

	defer func() {
		antsPool.ReleaseTimeout(5 * time.Second)
	}()

	array := []int{}
	for i := 0; i < 98; i++ {
		array = append(array, i)
	}

	fmt.Println(time.Now())
	antsPool.Submit(func() {
		time.Sleep(1 * time.Second)
		fmt.Println(">>>>>>222")
	})
	fmt.Println(time.Now())
	fmt.Println(">>>>>>next")
}

// 平台用户下WA账号登录信息 结构体  TaskUserWaLoginRecord
type TaskUserWaLoginRecord struct {
	Id      *int64 `json:"id" form:"id" gorm:"primarykey;column:id;comment:自增id;size:19;"`            //自增id
	UserId  *int64 `json:"userId" form:"userId" gorm:"column:user_id;comment:任务平台用户;size:19;"`        //任务平台用户
	WaId    *int64 `json:"waId" form:"waId" gorm:"column:wa_id;comment:平台用户绑定的 whatsapp 账号;size:19;"` //平台用户绑定的 whatsapp 账号
	StartTs *int   `json:"startTs" form:"startTs" gorm:"column:start_ts;comment:wa挂机开始时间;size:10;"`   //wa挂机开始时间
	EndTs   *int   `json:"endTs" form:"endTs" gorm:"column:end_ts;comment:wa挂机结束时间;size:10;"`         //wa挂机结束时间 - 默认值0，表示还未下线
}
