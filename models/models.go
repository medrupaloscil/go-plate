package models

import (
	"fmt"
	"go-plate/services"
	"reflect"
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Auto-migrate database models
func MigrateModels() bool {
	var models = []interface{}{
		&User{},
	}

	for _, s := range models {
		err := services.DB.AutoMigrate(s)
		if err != nil {
			fmt.Println("failed_automigrate", services.GetStructType(s), "model:", err.Error())
			return false
		}
	}

	return true
}

func UpdateElement(elem *reflect.Value, update *reflect.Value) {
	typeOfS := update.Type()

	for i := 0; i < update.NumField(); i++ {
		field := update.Field(i)
		if field.Kind() != reflect.Ptr {
			field = field.Addr()
		}
		if !field.IsNil() {
			if elem.FieldByName(typeOfS.Field(i).Name).CanSet() {
				switch field.Type().Elem().Kind() {
				case reflect.Float64:
					elem.FieldByName(typeOfS.Field(i).Name).SetFloat(field.Elem().Interface().(float64))
					break
				case reflect.Float32:
					elem.FieldByName(typeOfS.Field(i).Name).SetFloat(float64(field.Elem().Interface().(float32)))
					break
				case reflect.Int:
					elem.FieldByName(typeOfS.Field(i).Name).SetInt(field.Elem().Interface().(int64))
					break
				case reflect.Bool:
					elem.FieldByName(typeOfS.Field(i).Name).SetBool(field.Elem().Interface().(bool))
					break
				case reflect.String:
					elem.FieldByName(typeOfS.Field(i).Name).SetString(field.Elem().Interface().(string))
					break
				default:
					fmt.Printf("Cannot parse field %s\t of value:  %v\n", typeOfS.Field(i).Name, field.Elem().Interface())
					break
				}
			}
		}
	}
}
