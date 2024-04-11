package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Menu string
const Menu = `Databases Project Two Queries
    1. Add new employee
    2. View employee
    3. Modify employee
    4. Remove employee
    5. Add new dependent
    6. Remove dependent
    7. Add department
    8. View department
    9. Remove department 
    10. Add department location
    11. Remove department location
    `

func printMenu() {
	fmt.Print(Menu)
}

func handleUserInput() error {
	var userChoice uint8
	fmt.Println("Enter a number for one of the options above")
	fmt.Scan(&userChoice)

	switch userChoice {
	case 1:
		return addNewEmployee()
	case 2:
		return viewEmployee()
	case 3:
		return modifyEmployee()
	case 4:
		return removeEmployee()
	case 5:
		return addNewDependent()
	case 6:
		return removeDependent()
	case 7:
		return addNewDepartment()
	case 8:
		return viewDepartment()
	case 9:
		return removeDepartment()
	case 10:
		return addDepartmentLocation()
	case 11:
		return removeDepartmentLocation()
	default:
		return errors.New("invalid choice. Please enter a number that corresponds with a menu option")
	}
}

func addNewEmployee() error {

	var fName string
	var mInit string
	var lName string
	var ssn string
	var bDate string
	var address string
	var sex string
	var salary float64
	var superSsn string
	var dno uint

	fmt.Println("Enter employee information below")
	fmt.Println("first name")
	fmt.Scan(&fName)

	fmt.Println("middle initial")
	fmt.Scan(&mInit)

	fmt.Println("last name")
	fmt.Scan(&lName)

	fmt.Println("ssn")
	fmt.Scan(&ssn)

	fmt.Println("birth date")
	fmt.Scan(&bDate)
	date, timeErr := time.Parse(timeLayout, bDate)

	if timeErr != nil {
		panic("Wrong format, try \"year-month-day\" -> 2006-Jan-02 ")
	}

	fmt.Println("address")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		address = scanner.Text()
	}

	fmt.Println("sex")
	fmt.Scan(&sex)

	fmt.Println("salary")
	fmt.Scan(&salary)

	fmt.Println("super ssn")
	fmt.Scan(&superSsn)

	fmt.Println("department number")
	fmt.Scan(&dno)
	newEmployee := Employee{Fname: fName, Minit: mInit, Lname: lName, Ssn: ssn, Bdate: date, Address: address, Sex: sex, Salary: salary, SuperSsn: superSsn, Dno: dno}
	result := _edb.Create(&newEmployee)
	return result.Error
}

func viewEmployee() error {
	var emp Employee
	emp.find()
	emp.print()

	var supervisor Employee
	supervisor.Ssn = emp.SuperSsn
	_edb.First(&supervisor)
	fmt.Printf("Supervisor name: %s %s %s\n", supervisor.Fname, supervisor.Minit, supervisor.Lname)

	var dept Department
	_ddb.First(&dept, "Dnumber=?", emp.Dno)
	fmt.Printf("Department name: %s\n", dept.Dname)

	var dependents []Dependent
	result := _dedb.Where("Essn = ?", emp.Ssn).Find(&dependents)
	for _, dependent := range dependents {
		fmt.Printf("%+v\n", dependent)
	}

	return result.Error
}

func modifyEmployee() error {

	var emp Employee
	emp.find()
	emp.print()

	fmt.Println("Which attribute would you like to edit? (Enter the number)")
	fmt.Println(`1. Address
    2. Sex
    3. Salary
    4. SuperSsn
    5. Dno`)
	var userAnswer int
	fmt.Scan(&userAnswer)

	fmt.Println("What value do you want to replace with")
	switch userAnswer {
	case 1:
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			emp.Address = scanner.Text()
		}
	case 2:
		var userInput string
		fmt.Scan(&userInput)
		emp.Sex = userInput
	case 3:
		var userInput float64
		fmt.Scan(&userInput)
		emp.Salary = userInput
	case 4:
		var userInput string
		fmt.Scan(&userInput)
		emp.SuperSsn = userInput
	case 5:
		var userInput uint
		fmt.Scan(&userInput)
		emp.Dno = userInput
	default:
		panic("That was not an option")
	}

	result := _edb.Save(&emp)
	return result.Error
}

func removeEmployee() error {

	var emp Employee
	emp.find()
	emp.print()
	fmt.Println("Are you sure you want to delete (true or false)")
	var deleteEmp bool
	fmt.Scan(&deleteEmp)

	if deleteEmp {
		var dependent Dependent
		result := _dedb.First(&dependent, "Essn=?", emp.Ssn)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result = _edb.Delete(&dependent)
			return result.Error
		} else {
			fmt.Println("This employee has dependencies and you need to remove them before deleting")
		}
		return result.Error
	}
	return nil
}

func addNewDependent() error {
	var essn string
	var dependentName string
	var sex string
	var bDate string
	var relationship string

	fmt.Println("Enter Employee Ssn")
	fmt.Scan(&essn)

	// locking employee record
	var emp Employee
	_edb.Clauses(clause.Locking{Strength: "UPDATE"}).First(&emp, "Ssn = ?", essn)

	printAllDependents()

	fmt.Println("Enter the name of the dependent")
	fmt.Scan(&dependentName)

	fmt.Println("Enter the sex of the dependent")
	fmt.Scan(&sex)

	fmt.Println("Enter the birth date of the dependent")
	fmt.Scan(&bDate)
	date, timeErr := time.Parse(timeLayout, bDate)
	if timeErr != nil {
		panic("Wrong format, try \"year-month-day\" -> 2006-Jan-02 ")
	}

	fmt.Println("Enter the relationship to the employee")
	fmt.Scan(&relationship)

	newDependent := Dependent{Essn: essn, Dependent_Name: dependentName, Sex: sex, Bdate: date, Relationship: relationship}
	result := _dedb.Create(&newDependent)
	return result.Error

}

func removeDependent() error {
	var essn string
	var dName string

	fmt.Println("Enter Employee Ssn")
	fmt.Scan(&essn)

	// locking employee record
	var emp Employee
	_edb.Clauses(clause.Locking{Strength: "UPDATE"}).First(&emp, "Ssn = ?", essn)

	printAllDependents()

	fmt.Println("Enter dependent name that you want to delete")
	fmt.Scan(&dName)

	var dependent Dependent
	dependent.Essn = essn
	dependent.Dependent_Name = dName
	result := _dedb.Delete(&dependent)
	return result.Error
}

func addNewDepartment() error {
	var dName string
	var dNumber uint
	var mgrSsn string
	var mgrStartDate string

	fmt.Println("Enter department name")
	fmt.Scan(&dName)

	fmt.Println("Enter department number")
	fmt.Scan(&dNumber)

	fmt.Println("Enter manager ssn")
	fmt.Scan(&mgrSsn)

	fmt.Println("Enter manager start date")
	fmt.Scan(&mgrStartDate)
	date, timeErr := time.Parse(timeLayout, mgrStartDate)

	if timeErr != nil {
		panic("Wrong format, try \"year-month-day\" -> 2006-Jan-02 ")
	}

	newDepartment := Department{Dnumber: dNumber, Dname: dName, MgrSsn: mgrSsn, MgrStartDate: date}
	result := _ddb.Create(&newDepartment)
	return result.Error
}

func viewDepartment() error {
	var dNumber int
	fmt.Println("Enter department number")
	fmt.Scan(&dNumber)

	printAllDepartments()

	var dep Department
	result := _ddb.First(&dep, "Dnumber = ?", dNumber)

	var emp Employee
	_edb.First(&emp, "Ssn = ?", dep.MgrSsn)
	fmt.Printf("Manager name: %s\n", emp.Fname)

	printAllDeptLocations()

	return result.Error
}

func removeDepartment() error {
	var dNumber int
	fmt.Println("Enter department number")
	fmt.Scan(&dNumber)

	var dep Department
	_ddb.Clauses(clause.Locking{Strength: "UPDATE"}).First(&dep, "Dnumber = ?", dNumber)
	fmt.Printf("%+v\n", dep)

	var delDepartment bool
	fmt.Println("Are you sure you want to delete the department (true or false)")
	fmt.Scan(&delDepartment)

	if delDepartment {
		var deptLoc DeptLocation
		result := _dldb.First(&deptLoc, "Dnumber=?", dep.Dnumber)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			_ddb.Delete(&dep)
		} else {
			fmt.Println("This employee has dependencies and you need to remove them before deleting")
		}
	}
	return nil
}

func addDepartmentLocation() error {
	var dNumber uint
	var dLocation string
	fmt.Println("Enter department number")
	fmt.Scan(&dNumber)
	var dep Department

	_ddb.Clauses(clause.Locking{Strength: "UPDATE"}).First(&dep, "Dnumber = ?", dNumber)

	printAllDeptLocations()

	fmt.Println("Enter department location")
	fmt.Scan(&dLocation)

	newDepLoc := DeptLocation{dNumber, dLocation}
	result := _dldb.Create(&newDepLoc)

	return result.Error
}

func removeDepartmentLocation() error {
	var dNumber int
	var dLocation string
	var deptLoc DeptLocation
	var dep Department
	fmt.Println("Enter department number")
	fmt.Scan(&dNumber)

	_ddb.Clauses(clause.Locking{Strength: "UPDATE"}).First(&dep, "Dnumber = ?", dNumber)

	printAllDeptLocations()

	fmt.Println("Enter department location")
	fmt.Scan(&dLocation)

	result := _dldb.Delete(&deptLoc, "Dnumber = ? AND Dlocation = ?", dNumber, dLocation)
	return result.Error
}
