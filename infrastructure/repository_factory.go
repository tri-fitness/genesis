package infrastructure

import work "github.com/freerware/work"

type RepositoryFactory interface {
	Account(work.Unit) AccountRepository
}

type repositoryFactory struct {
	queryFactory QueryFactory
}

func NewRepositoryFactory(queryFactory QueryFactory) RepositoryFactory {
	return &repositoryFactory{
		queryFactory: queryFactory,
	}
}

func (r *repositoryFactory) Account(unit work.Unit) AccountRepository {
	return NewAccountRepository(unit, r.queryFactory)
}
