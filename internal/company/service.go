package company

type Repository interface {
	FindById(id int) (*Company, error)
	FindAll() (*[]Company, error)
	DeleteById(id int) error
	Insert(*Company) (*Company, error)
	Update(*Company) (*Company, error)
	Exists(id int) (bool, error)
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{r}
}

func (s *Service) GetCompany(id int) (*Company, error) {
	c, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Service) GetCompanies() (*[]Company, error) {
	cList, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return cList, nil
}

func (s *Service) CreateCompany(c *Company) (*Company, error) {
	newComp, err := s.repo.Insert(c)
	if err != nil {
		return nil, err
	}

	return newComp, nil
}

func (s *Service) UpdateCompany(id int, newData *Company) (*Company, error) {
	cmp, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	cmp.Name = newData.Name

	updatedCmp, err := s.repo.Update(cmp)
	if err != nil {
		return nil, err
	}

	return updatedCmp, nil
}

func (s *Service) DeleteCompany(id int) error {
	exists, err := s.repo.Exists(id)
	if err != nil {
		return err
	}

	if !exists {
		return ErrCompanyNotFound
	}

	return s.repo.DeleteById(id)
}
