package tests

import (
	"testing"

	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/sqlite"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	queries map[string]*TestQuery
}

func TestDBTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &DBTestSuite{
		queries: make(map[string]*TestQuery),
	})
}

func (s *DBTestSuite) SetupSuite() {
	s.queries = NewTestQueryBuilder().All("", false)
	for _, query := range s.queries {
		query.CreateTable(TestTableProducts)
	}
}

func (s *DBTestSuite) TearDownSuite() {
	if s.queries[sqlite.Name] != nil {
		docker, err := s.queries[sqlite.Name].Driver().Docker()
		s.NoError(err)
		s.NoError(docker.Shutdown())
	}
}

func (s *DBTestSuite) TestInsert_First_Get() {
	for driver, query := range s.queries {
		now := carbon.NewDateTime(carbon.FromDateTime(2025, 1, 2, 3, 4, 5))

		s.Run(driver, func() {
			s.Run("single struct", func() {
				result, err := query.DB().Table("products").Insert(Product{
					Name: "single struct",
					Model: Model{
						Timestamps: Timestamps{
							CreatedAt: now,
							UpdatedAt: now,
						},
					},
				})

				s.NoError(err)
				s.Equal(int64(1), result.RowsAffected)

				var product Product
				err = query.DB().Table("products").Where("name", "single struct").Where("deleted_at", nil).First(&product)
				s.NoError(err)
				s.True(product.ID > 0)
				s.Equal("single struct", product.Name)
				s.Equal(now, product.CreatedAt)
				s.Equal(now, product.UpdatedAt)
				s.False(product.DeletedAt.Valid)
			})

			s.Run("multiple structs", func() {
				result, err := query.DB().Table("products").Insert([]Product{
					{
						Name: "multiple structs1",
						Model: Model{
							Timestamps: Timestamps{
								CreatedAt: now,
								UpdatedAt: now,
							},
						},
					},
					{
						Name: "multiple structs2",
					},
				})
				s.NoError(err)
				s.Equal(int64(2), result.RowsAffected)

				var products []Product
				err = query.DB().Table("products").Where("name", []string{"multiple structs1", "multiple structs2"}).Where("deleted_at", nil).Get(&products)
				s.NoError(err)
				s.Equal(2, len(products))
				s.Equal("multiple structs1", products[0].Name)
				s.Equal("multiple structs2", products[1].Name)
			})

			s.Run("single map", func() {
				result, err := query.DB().Table("products").Insert(map[string]any{
					"name":       "single map",
					"created_at": now,
					"updated_at": &now,
				})
				s.NoError(err)
				s.Equal(int64(1), result.RowsAffected)

				var product Product
				err = query.DB().Table("products").Where("name", "single map").Where("deleted_at", nil).First(&product)
				s.NoError(err)
				s.Equal("single map", product.Name)
				s.Equal(now, product.CreatedAt)
				s.Equal(now, product.UpdatedAt)
				s.False(product.DeletedAt.Valid)
			})

			s.Run("multiple map", func() {
				result, err := query.DB().Table("products").Insert([]map[string]any{
					{
						"name":       "multiple map1",
						"created_at": now,
						"updated_at": &now,
					},
					{
						"name": "multiple map2",
					},
				})
				s.NoError(err)
				s.Equal(int64(2), result.RowsAffected)

				var products []Product
				err = query.DB().Table("products").Where("name", []string{"multiple map1", "multiple map2"}).Where("deleted_at", nil).Get(&products)
				s.NoError(err)
				s.Equal(2, len(products))
				s.Equal("multiple map1", products[0].Name)
				s.Equal("multiple map2", products[1].Name)
			})
		})
	}
}

func (s *DBTestSuite) TestWhere() {
	for driver, query := range s.queries {
		s.Run(driver, func() {
			now := carbon.NewDateTime(carbon.FromDateTime(2025, 1, 2, 3, 4, 5))
			query.DB().Table("products").Insert(Product{
				Name: "where model",
				Model: Model{
					Timestamps: Timestamps{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			})

			s.Run("simple where condition", func() {
				var product Product
				err := query.DB().Table("products").Where("name", "where model").First(&product)
				s.NoError(err)
				s.Equal("where model", product.Name)
			})

			s.Run("where with multiple arguments", func() {
				var products []Product
				err := query.DB().Table("products").Where("name", []string{"where model", "where model1"}).Get(&products)
				s.NoError(err)
				s.Equal(1, len(products))
				s.Equal("where model", products[0].Name)
			})

			s.Run("where with raw query", func() {
				var product Product
				err := query.DB().Table("products").Where("name = ?", "where model").First(&product)
				s.NoError(err)
				s.Equal("where model", product.Name)
			})
		})
	}
}
