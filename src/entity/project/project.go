package project

import (
	"errors"

	"github.com/gofrs/uuid"
)

type Project struct {
	ID      uuid.UUID
	Name    string
	OwnerID uuid.UUID
}

func New(opts ...ProjectOption) (*Project, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	project := &Project{
		ID:   uuid,
		Name: "unnamed",
	}

	errs := []error{}
	for _, opt := range opts {
		errs = append(errs, opt(project))
	}
	err = errors.Join(errs...)
	if err != nil {
		return nil, err
	}

	return project, nil
}

type ProjectOption func(p *Project) error

func WithName(name string) ProjectOption {
	return func(p *Project) error {
		p.Name = name
		return nil
	}
}

func WithOwnerID(ownerID uuid.UUID) ProjectOption {
	return func(p *Project) error {
		p.OwnerID = ownerID
		return nil
	}
}
