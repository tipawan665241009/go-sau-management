package main

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func studentRoute(app *fiber.App, db *gorm.DB) {

	// POST: Create student
	app.Post("/api/students", func(c fiber.Ctx) error {
		student := new(Student)

		if err := c.Bind().Body(student); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		if student.StuName == "" ||
			student.StuLevel < 1 ||
			student.StuGPA < 0 ||
			student.StuGPA > 4 {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid student data",
			})
		}

		if err := createStudent(db, student); err != nil {
			log.Printf("Create error: %v", err)

			return c.Status(500).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"message": "Student created successfully",
			"data":    student,
		})
	})

	// GET: All students / Search by name or level
	app.Get("/api/students", func(c fiber.Ctx) error {
		stuName := strings.TrimSpace(c.Query("StuName"))
		stuLevel := strings.TrimSpace(c.Query("StuLevel"))

		var students []Student
		var err error

		// ค้นหานักศึกษาด้วยชื่อ
		if stuName != "" {
			students, err = getStudentByName(db, stuName)

			// ค้นหานักศึกษาด้วยระดับชั้น
		} else if stuLevel != "" {
			level, parseErr := strconv.Atoi(stuLevel)

			if parseErr != nil || level < 1 {
				return c.Status(400).JSON(fiber.Map{
					"message": "Invalid student level",
				})
			}

			students, err = getStudentByLevel(db, level)

			// ดึงข้อมูลนักศึกษาทั้งหมด
		} else {
			students, err = getAllStudent(db)
		}

		// ตรวจสอบ Error จาก Database
		if err != nil {
			log.Printf("Get students error: %v", err)

			return c.Status(500).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		if students == nil {
			students = []Student{}
		}

		return c.Status(200).JSON(fiber.Map{
			"message": "Success",
			"data":    students,
		})
	})

	// GET: Student by ID
	app.Get("/api/students/:id", func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		if err != nil || id <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request parameter",
			})
		}

		student, err := getStudentByID(db, id)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"message": "Student not found",
			})
		}

		if err != nil {
			log.Printf("Get student by ID error: %v", err)

			return c.Status(500).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"message": "Success",
			"data":    student,
		})
	})

	// API Put for Update (Endpoint:/api/students/:id)
	app.Put("/api/students/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		idInt, error := strconv.Atoi(id)
		if error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request parameter",
			})
		}

		student := new(Student)

		if error := c.Bind().Body(student); error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		if error := updateStudent(db, idInt, student); error != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"message": "Student updated successfully",
			"data":    student,
		})
	})

	// DELETE: Delete student by ID
	app.Delete("/api/students/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		idInt, error := strconv.Atoi(id)
		if error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request parameter",
			})
		}

		if error := deleteStudent(db, idInt); error != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"message": "Student delete successfully",
		})
	})
}
