package student

import (
	"context"

	student_pb "go_test/api/student"
	"go_test/pkg/ent"
	"go_test/pkg/ent/student"

	"github.com/sirupsen/logrus"
)

type service struct {
	student_pb.UnimplementedStudentSvcServer
	logger *logrus.Logger
	DB     *ent.Client
}

func New(logger *logrus.Logger, db *ent.Client) *service {
	return &service{
		logger: logger,
		DB:     db,
	}
}

func (s *service) GetStudent(
	ctx context.Context,
	req *student_pb.GetStudentRequest,
) (*student_pb.GetStudentResponse, error) {
	stu, err := s.DB.Student.Query().WithTeacher().Where(student.StuID(req.StuId)).First(ctx)
	if err != nil {
		return nil, err
	}
	var teacherName string
	if stu.Edges.Teacher != nil {
		teacherName = stu.Edges.Teacher.TeacherName
	}
	return &student_pb.GetStudentResponse{
		Status: 0,
		Msg:    "OK",
		Data: &student_pb.StudentWithTeacher{
			Student: &student_pb.Student{
				StuId: stu.StuID,
				Name:  stu.Name,
				Age:   stu.Age,
			},
			Teacher: teacherName,
		},
	}, nil
}

func (s *service) CreateStudent(
	ctx context.Context,
	req *student_pb.Student,
) (*student_pb.CreateStudentResponse, error) {
	stu, err := s.DB.Student.Create().
		SetStuID(req.StuId).
		SetAge(req.Age).
		SetName(req.Name).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &student_pb.CreateStudentResponse{
		Status: 0,
		Msg:    "OK",
		Data: &student_pb.Student{
			StuId: stu.StuID,
			Name:  stu.Name,
			Age:   stu.Age,
		},
	}, err
}
