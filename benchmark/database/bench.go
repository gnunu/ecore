package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Record struct {
	Name string `gorm:"primaryKey"`
	Tag  string `gorm:"not null"`
	Seq  int    `gorm:"not null"`
}

type database struct {
	db *gorm.DB
}

func (d *database) build(name string) {
	ds := fmt.Sprintf("user=bench password=bench dbname=%s host=postgres port=5432 sslmode=disable", name)
	db, err := gorm.Open(postgres.Open(ds), &gorm.Config{})
	if err != nil {
		panic("failed to perform migrations: " + err.Error())
	}

	err = db.AutoMigrate(&Record{})
	if err != nil {
		panic("failed to perform migrations: " + err.Error())
	}
	d.db = db
}

func (d *database) clear() {
	var rs []Record
	d.db.Find(&rs)

	for _, r := range rs {
		d.db.Delete(&r)
	}
}

func (d *database) create(i int) {
	r := Record{
		Name: fmt.Sprintf("foo{%d}", i),
		Tag:  fmt.Sprintf("bar=%d=", i),
		Seq:  i,
	}

	result := d.db.Create(&r)
	if result.Error != nil {
		panic("failed to create record: " + result.Error.Error())
	}
	fmt.Printf("record %s %s %d was created\n", r.Name, r.Tag, r.Seq)
}

func (d *database) read(i int) {
	var rs []Record

	result := d.db.Where("Name = ?", fmt.Sprintf("foo{%d}", i)).Where("Tag = ?", fmt.Sprintf("bar=%d=", i)).Find(&rs)
	if result.Error != nil {
		panic("failed to retrieve record: " + result.Error.Error())
	}

	for _, r := range rs {
		fmt.Printf("Name: %s, Tag: %s, Seq: %d\n", r.Name, r.Tag, r.Seq)
	}
}

func (d *database) update(i int) {
	var rs []Record

	result := d.db.Where("Name = ?", fmt.Sprintf("foo{%d}", i)).Find(&rs)
	if result.Error != nil {
		panic("failed to retrieve record: " + result.Error.Error())
	}

	for _, r := range rs {
		r.Seq += 100
		result = d.db.Save(&r)
		if result.Error != nil {
			panic("failed to update user: " + result.Error.Error())
		}
		fmt.Printf("record retrieved: %s %s %d\n", r.Name, r.Tag, r.Seq)
	}
}

func (d *database) delete(i int) {
	var rs []Record

	result := d.db.Where("Name = ?", fmt.Sprintf("foo{%d}", i)).Find(&rs)
	if result.Error != nil {
		panic("failed to retrieve record: " + result.Error.Error())
	}

	for _, r := range rs {
		d.db.Delete(&r)
		if result.Error != nil {
			panic("failed to delete record: " + result.Error.Error())
		} else if result.RowsAffected == 0 {
			panic("no record was deleted")
		} else {
			fmt.Println("record deleted successfully")
		}
	}
}

func (d *database) bench() {
	for k := 0; ; k++ {
		for i := k; i < k+10; i++ {
			d.create(i)
		}
		for i := k; i < k+10; i++ {
			d.read(i)
		}
		for i := k; i < k+10; i++ {
			d.update(i)
		}
		for i := k; i < k+10; i++ {
			d.delete(i)
		}
	}
}
