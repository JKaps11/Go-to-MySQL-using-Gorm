package main

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _db *gorm.DB

var _edb *gorm.DB
var _ddb *gorm.DB
var _dldb *gorm.DB
var _dedb *gorm.DB

const timeLayout = "2006-Jan-02"

func connectToDatabase(password string) {

	dsn := "root:" + password + "@tcp(127.0.0.1:3306)/COMPANY?charset=utf8&parseTime=True&loc=Local"

	gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	_db = gormdb
	if err != nil {
		panic("failed to make mysql db to gorm db")
	}
	initalizeTables()
}

func initalizeTables() {
	_edb = _db.Table("EMPLOYEE")
	_ddb = _db.Table("DEPARTMENT")
	_dldb = _db.Table("DEPT_LOCATIONS")
	_dedb = _db.Table("DEPENDENT")
}

// BeforeCreate method activates before creation attempt
func (emp *Employee) BeforeCreate(*gorm.DB) error {

	if len(emp.Fname) > 15 {
		return errors.New("first name to big. needs to be less than 16 characters")
	}

	if len(emp.Minit) > 1 {
		return errors.New("middle initial can only be one character")
	}

	if len(emp.Lname) > 15 {
		return errors.New("last name to big, needs to be less than 16 characters")
	}

	if len(emp.Ssn) > 9 {
		return errors.New("ssn to big, needs to be less than 10 characters")
	}

	if len(emp.Address) > 30 {
		return errors.New("address to big, needs to be less than 31 characters")
	}

	if len(emp.Sex) > 1 {
		return errors.New("sex can only be one character")
	}

	// need to implement check for salary length

	if len(emp.SuperSsn) > 9 {
		return errors.New("super ssn to big, needs to be less than 10 characters")
	}

	return nil
}

func (emp *Employee) find() {
	fmt.Println("Enter a Ssn")
	var ssn string
	fmt.Scan(&ssn)

	// locks for entire transaction
	result := _edb.Clauses(clause.Locking{Strength: "UPDATE"}).First(&emp, "Ssn=?", ssn) // SELECT * FROM COMPANY.EMPLOYEE WHERE Ssn="ssn"
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		panic("Employee with ssn: " + ssn + " not found")
	}
}

func (emp *Employee) print() {
	fmt.Printf("%+v\n", emp)
}

// BeforeCreate method activates before creation attempt
func (dependent *Dependent) BeforeCreate(*gorm.DB) error {
	var emp Employee
	if errors.Is(_edb.First(&emp, "Ssn=?", dependent.Essn).Error, gorm.ErrRecordNotFound) {
		return errors.New("no employee found with the provided essn")
	}

	if len(dependent.Dependent_Name) > 15 {
		return errors.New("Dependent name to big, needs to be less than 16 characters")
	}

	if len(dependent.Sex) > 1 {
		return errors.New("sex should only be a single character")
	}

	if len(dependent.Relationship) > 8 {
		return errors.New("string size to big, need to be less than 10 characters")
	}
	return nil
}

func printAllDependents() {
	var dependents []Dependent
	result := _dedb.Find(&dependents)
	if result.Error != nil {
		panic("Problem with getting dependents from database")
	}

	for _, dependent := range dependents {
		fmt.Printf("%+v\n", dependent)
	}
}

func printAllDepartments() {
	var departments []Department
	result := _ddb.Find(&departments)
	if result.Error != nil {
		panic("Problem with getting departments from databse")
	}

	for _, department := range departments {
		fmt.Printf("%+v\n,", department)
	}
}

// BeforeCreate method activates before creation attempt
func (department *Department) BeforeCreate(*gorm.DB) error {
	if len(department.Dname) > 15 {
		return errors.New("department name too big, needs to be less than 16 characters")
	}

	var emp Employee
	if errors.Is(_edb.First(&emp, "Ssn=?", department.MgrSsn).Error, gorm.ErrRecordNotFound) {
		return errors.New("no employee found with the provided with the mgr_ssn")
	}
	return nil
}

// BeforeCreate method activates before creation attempt
func (deptLoc *DeptLocation) BeforeCreate(*gorm.DB) error {
	if len(deptLoc.Dlocation) > 15 {
		return errors.New("department location to big, needs to be less than 16 characters")
	}

	var dep Department
	if errors.Is(_ddb.First(&dep, "Dnumber=?", deptLoc.Dnumber).Error, gorm.ErrRecordNotFound) {
		return errors.New("no department found with the provided dnumber")
	}
	return nil
}

func printAllDeptLocations() {
	var deptLocs []DeptLocation
	result := _dldb.Find(&deptLocs)
	if result.Error != nil {
		panic("Problem with getting department locations from database")
	}

	for _, deptLoc := range deptLocs {
		fmt.Printf("%+v\n", deptLoc)
	}
}
