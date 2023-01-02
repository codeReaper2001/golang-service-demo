// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"go_test/pkg/ent/student"
	"go_test/pkg/ent/teacher"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// TeacherCreate is the builder for creating a Teacher entity.
type TeacherCreate struct {
	config
	mutation *TeacherMutation
	hooks    []Hook
}

// SetTeacherID sets the "teacher_id" field.
func (tc *TeacherCreate) SetTeacherID(i int32) *TeacherCreate {
	tc.mutation.SetTeacherID(i)
	return tc
}

// SetTeacherName sets the "teacher_name" field.
func (tc *TeacherCreate) SetTeacherName(s string) *TeacherCreate {
	tc.mutation.SetTeacherName(s)
	return tc
}

// AddStudentIDs adds the "students" edge to the Student entity by IDs.
func (tc *TeacherCreate) AddStudentIDs(ids ...int) *TeacherCreate {
	tc.mutation.AddStudentIDs(ids...)
	return tc
}

// AddStudents adds the "students" edges to the Student entity.
func (tc *TeacherCreate) AddStudents(s ...*Student) *TeacherCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return tc.AddStudentIDs(ids...)
}

// Mutation returns the TeacherMutation object of the builder.
func (tc *TeacherCreate) Mutation() *TeacherMutation {
	return tc.mutation
}

// Save creates the Teacher in the database.
func (tc *TeacherCreate) Save(ctx context.Context) (*Teacher, error) {
	var (
		err  error
		node *Teacher
	)
	if len(tc.hooks) == 0 {
		if err = tc.check(); err != nil {
			return nil, err
		}
		node, err = tc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TeacherMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = tc.check(); err != nil {
				return nil, err
			}
			tc.mutation = mutation
			if node, err = tc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(tc.hooks) - 1; i >= 0; i-- {
			if tc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = tc.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, tc.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Teacher)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from TeacherMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TeacherCreate) SaveX(ctx context.Context) *Teacher {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TeacherCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TeacherCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tc *TeacherCreate) check() error {
	if _, ok := tc.mutation.TeacherID(); !ok {
		return &ValidationError{Name: "teacher_id", err: errors.New(`ent: missing required field "Teacher.teacher_id"`)}
	}
	if _, ok := tc.mutation.TeacherName(); !ok {
		return &ValidationError{Name: "teacher_name", err: errors.New(`ent: missing required field "Teacher.teacher_name"`)}
	}
	return nil
}

func (tc *TeacherCreate) sqlSave(ctx context.Context) (*Teacher, error) {
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (tc *TeacherCreate) createSpec() (*Teacher, *sqlgraph.CreateSpec) {
	var (
		_node = &Teacher{config: tc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: teacher.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: teacher.FieldID,
			},
		}
	)
	if value, ok := tc.mutation.TeacherID(); ok {
		_spec.SetField(teacher.FieldTeacherID, field.TypeInt32, value)
		_node.TeacherID = value
	}
	if value, ok := tc.mutation.TeacherName(); ok {
		_spec.SetField(teacher.FieldTeacherName, field.TypeString, value)
		_node.TeacherName = value
	}
	if nodes := tc.mutation.StudentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   teacher.StudentsTable,
			Columns: []string{teacher.StudentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: student.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TeacherCreateBulk is the builder for creating many Teacher entities in bulk.
type TeacherCreateBulk struct {
	config
	builders []*TeacherCreate
}

// Save creates the Teacher entities in the database.
func (tcb *TeacherCreateBulk) Save(ctx context.Context) ([]*Teacher, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Teacher, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TeacherMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TeacherCreateBulk) SaveX(ctx context.Context) []*Teacher {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TeacherCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TeacherCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}