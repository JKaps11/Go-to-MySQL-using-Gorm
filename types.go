package main

import (
	"time"
)

// Employee Enitity
type Employee struct {
	Fname    string
	Minit    string
	Lname    string
	Ssn      string `gorm:"primaryKey"`
	Bdate    time.Time
	Address  string
	Sex      string
	Salary   float64
	SuperSsn string
	Dno      uint
}

// Department Entity
type Department struct {
	Dnumber      uint `gorm:"primaryKey"`
	Dname        string
	MgrSsn       string    `gorm:"column:Mgr_ssn"`
	MgrStartDate time.Time `gorm:"column:Mgr_start_date"`
}

// DeptLocation Entity
type DeptLocation struct {
	Dnumber   uint `gorm:"primaryKey"`
	Dlocation string
}

// Dependent Entity
type Dependent struct {
	Essn           string `gorm:"primaryKey"`
	Dependent_Name string `gorm:"primaryKey; column:Dependent_name"`
	Sex            string
	Bdate          time.Time
	Relationship   string
}
