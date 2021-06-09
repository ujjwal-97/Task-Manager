package models

import (
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" `
	Title            string             `json:"title" bson:"title"`
	Description      string             `json:"description" bson:"description" `
	Deadline         primitive.DateTime `json:"deadline,omitempty" bson:"deadline,omitempty"`
	CreatedAt        primitive.DateTime `json:"createdAt" bson:"createdAt"`
	TaskUser         *User              `json:"user,omitempty" bson:"user,omitempty"`
	Status           string             `json:"status" bson:"status"`
	CronID           cron.EntryID       `json:"-" bson:"cronid"`
	SnapshotSchedule *ScheduleSnapshot  `json:"snapshotSchedule" bson:"snapshotSchedule"`
}

type ScheduleSnapshot struct {
	Periodic bool      `json:"periodic" bson:"periodic"`
	Schedule *Schedule `json:"schedule" bson:"schedule"`
	Interval *Interval `json:"interval" bson:"interval"`
}
type Interval struct {
	Minute int `json:"minute" bson:"minute"`
	Hour   int `json:"hour" bson:"hour"`
	Day    int `json:"day" bson:"day"`
	Month  int `json:"month" bson:"month"`
}
type Schedule struct {
	Minute  int `json:"minute" bson:"minute"`
	Hour    int `json:"hour" bson:"hour"`
	Day     int `json:"day" bson:"day"`
	Month   int `json:"month" bson:"month"`
	Weekday int `json:"weekday" bson:"weekday"`
}
