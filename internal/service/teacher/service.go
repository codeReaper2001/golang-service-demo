package teacher

import (
	"context"
	teacher_pb "go_test/api/teacher"
	"go_test/pkg/ent"
	"go_test/pkg/ent/student"
	"go_test/pkg/ent/teacher"

	"entgo.io/ent/dialect/sql"
	"github.com/sirupsen/logrus"
)

type service struct {
	logger *logrus.Logger
	DB     *ent.Client
}

func New(logger *logrus.Logger, db *ent.Client) *service {
	return &service{
		logger: logger,
		DB:     db,
	}
}

func (s *service) GetTeacher(
	ctx context.Context,
	req *teacher_pb.GetTeacherRequest,
) (*teacher_pb.GetTeacherResponse, error) {
	teacher, err := s.DB.Teacher.Query().WithStudents(func(sq *ent.StudentQuery) {
		t := sql.Table(student.Table)
		sq.Order(func(s *sql.Selector) {
			s.OrderBy(t.C(student.FieldStuID))
		})
	}).Where(teacher.TeacherID(req.TeacherId)).First(ctx)
	if ent.IsNotFound(err) {
		return &teacher_pb.GetTeacherResponse{
			Status: 404,
			Msg:    "teacher not found",
			Data:   &teacher_pb.TeacherWithStudents{},
		}, nil
	}
	if err != nil {
		return nil, err
	}
	stus := teacher.Edges.Students
	stuRes := []*teacher_pb.Student{}
	for _, stu := range stus {
		stuRes = append(stuRes, &teacher_pb.Student{
			StuId: stu.StuID,
			Name:  stu.Name,
			Age:   stu.Age,
		})
	}
	return &teacher_pb.GetTeacherResponse{
		Status: 0,
		Msg:    "OK",
		Data: &teacher_pb.TeacherWithStudents{
			TeacherId: teacher.TeacherID,
			Name:      teacher.TeacherName,
			Students:  stuRes,
		},
	}, nil
}
