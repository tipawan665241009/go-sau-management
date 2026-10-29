package main

import "gorm.io/gorm"

// Create student
func createStudent(db *gorm.DB, student *Student) error {
	return db.Create(student).Error
}

// Get all students
func getAllStudent(db *gorm.DB) ([]Student, error) {
	var students []Student

	result := db.Find(&students)

	return students, result.Error
}

// Get student by ID
func getStudentByID(db *gorm.DB, id int) (*Student, error) {
	var student Student

	result := db.First(&student, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &student, nil
}

// Search students by name
func getStudentByName(db *gorm.DB, name string) ([]Student, error) {
	var students []Student

	result := db.Where(
		"stu_name LIKE ?",
		"%"+name+"%",
	).Find(&students)

	return students, result.Error
}

// Search students by level
func getStudentByLevel(db *gorm.DB, level int) ([]Student, error) {
	var students []Student

	result := db.Where(
		"stu_level = ?",
		level,
	).Find(&students)

	return students, result.Error
}

func updateStudent(db *gorm.DB, id int, student *Student) error {
	return db.Model(&Student{}).Where("id = ?", id).Updates(student).Error
}

func deleteStudent(db *gorm.DB, id int) error {
	return db.Delete(&Student{}, id).Error
}
