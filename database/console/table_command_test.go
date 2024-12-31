package console

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/database/schema"
	mocksconfig "github.com/goravel/framework/mocks/config"
	mocksconsole "github.com/goravel/framework/mocks/console"
	mocksschema "github.com/goravel/framework/mocks/database/schema"
)

func TestTableCommand(t *testing.T) {
	var (
		mockContext *mocksconsole.Context
		mockConfig  *mocksconfig.Config
		mockSchema  *mocksschema.Schema
	)
	beforeEach := func() {
		mockContext = mocksconsole.NewContext(t)
		mockConfig = mocksconfig.NewConfig(t)
		mockSchema = mocksschema.NewSchema(t)
	}
	successCaseExpected := [][2]string{
		{"<fg=green;op=bold>test</>", ""},
		{"Columns", "1"},
		{"Size", "0.000MiB"},
		{"<fg=green;op=bold>Column</>", "Type"},
		{"foo <fg=gray>autoincrement, int, nullable</>", "<fg=gray>bar</> int"},
		{"<fg=green;op=bold>Index</>", ""},
		{"index_foo <fg=gray>foo, bar</>", "compound, unique, primary"},
		{"<fg=green;op=bold>Foreign Key</>", "On Update / On Delete"},
		{"fk_foo <fg=gray>foo references baz on bar</>", "restrict / cascade"},
	}
	tests := []struct {
		name  string
		setup func()
	}{
		{
			name: "get tables failed",
			setup: func() {
				mockContext.EXPECT().NewLine().Once()
				mockContext.EXPECT().Option("database").Return("").Once()
				mockSchema.EXPECT().Connection("").Return(mockSchema).Once()
				mockContext.EXPECT().Argument(0).Return("").Once()
				mockSchema.EXPECT().GetTables().Return(nil, assert.AnError).Once()
				mockContext.EXPECT().Error(fmt.Sprintf("Failed to get tables: %s", assert.AnError.Error())).Once()
			},
		},
		{
			name: "table not found",
			setup: func() {
				mockContext.EXPECT().NewLine().Times(2)
				mockContext.EXPECT().Option("database").Return("test").Once()
				mockSchema.EXPECT().Connection("test").Return(mockSchema).Once()
				mockContext.EXPECT().Argument(0).Return("test").Once()
				mockSchema.EXPECT().GetTables().Return(nil, nil).Once()
				mockContext.EXPECT().Warning("Table 'test' doesn't exist.").Once()
			},
		},
		{
			name: "choice table canceled",
			setup: func() {
				mockContext.EXPECT().NewLine().Times(1)
				mockContext.EXPECT().Option("database").Return("test").Once()
				mockSchema.EXPECT().Connection("test").Return(mockSchema).Once()
				mockContext.EXPECT().Argument(0).Return("").Once()
				mockSchema.EXPECT().GetTables().Return(nil, nil).Once()
				mockContext.EXPECT().Choice("Which table would you like to inspect?",
					[]console.Choice(nil)).Return("", assert.AnError).Once()
				mockContext.EXPECT().Line(assert.AnError.Error()).Once()
			},
		},
		{
			name: "get columns failed",
			setup: func() {
				mockContext.EXPECT().NewLine().Once()
				mockContext.EXPECT().Option("database").Return("test").Once()
				mockSchema.EXPECT().Connection("test").Return(mockSchema).Once()
				mockContext.EXPECT().Argument(0).Return("").Once()
				mockSchema.EXPECT().GetTables().Return([]schema.Table{{Name: "test"}}, nil).Once()
				mockContext.EXPECT().Choice("Which table would you like to inspect?",
					[]console.Choice{{Key: "test", Value: "test"}}).Return("test", nil).Once()
				mockSchema.EXPECT().GetColumns("test").Return(nil, assert.AnError).Once()
				mockContext.EXPECT().Error(fmt.Sprintf("Failed to get columns: %s", assert.AnError.Error())).Once()
			},
		},
		{
			name: "get indexes failed",
			setup: func() {
				mockContext.EXPECT().NewLine().Once()
				mockContext.EXPECT().Option("database").Return("test").Once()
				mockSchema.EXPECT().Connection("test").Return(mockSchema).Once()
				mockContext.EXPECT().Argument(0).Return("").Once()
				mockSchema.EXPECT().GetTables().Return([]schema.Table{{Name: "test"}}, nil).Once()
				mockContext.EXPECT().Choice("Which table would you like to inspect?",
					[]console.Choice{{Key: "test", Value: "test"}}).Return("test", nil).Once()
				mockSchema.EXPECT().GetColumns("test").Return(nil, nil).Once()
				mockSchema.EXPECT().GetIndexes("test").Return(nil, assert.AnError).Once()
				mockContext.EXPECT().Error(fmt.Sprintf("Failed to get indexes: %s", assert.AnError.Error())).Once()
			},
		},
		{
			name: "get foreign keys failed",
			setup: func() {
				mockContext.EXPECT().NewLine().Once()
				mockContext.EXPECT().Option("database").Return("test").Once()
				mockSchema.EXPECT().Connection("test").Return(mockSchema).Once()
				mockContext.EXPECT().Argument(0).Return("").Once()
				mockSchema.EXPECT().GetTables().Return([]schema.Table{{Name: "test"}}, nil).Once()
				mockContext.EXPECT().Choice("Which table would you like to inspect?",
					[]console.Choice{{Key: "test", Value: "test"}}).Return("test", nil).Once()
				mockSchema.EXPECT().GetColumns("test").Return(nil, nil).Once()
				mockSchema.EXPECT().GetIndexes("test").Return(nil, nil).Once()
				mockSchema.EXPECT().GetForeignKeys("test").Return(nil, assert.AnError).Once()
				mockContext.EXPECT().Error(fmt.Sprintf("Failed to get foreign keys: %s", assert.AnError.Error())).Once()
			},
		},
		{
			name: "success",
			setup: func() {
				mockContext.EXPECT().NewLine().Times(5)
				mockContext.EXPECT().Option("database").Return("test").Once()
				mockSchema.EXPECT().Connection("test").Return(mockSchema).Once()
				mockContext.EXPECT().Argument(0).Return("").Once()
				mockSchema.EXPECT().GetTables().Return([]schema.Table{{Name: "test"}}, nil).Once()
				mockContext.EXPECT().Choice("Which table would you like to inspect?",
					[]console.Choice{{Key: "test", Value: "test"}}).Return("test", nil).Once()
				mockSchema.EXPECT().GetColumns("test").Return([]schema.Column{
					{Name: "foo", Type: "int", TypeName: "int", Autoincrement: true, Nullable: true, Default: "bar"},
				}, nil).Once()
				mockSchema.EXPECT().GetIndexes("test").Return([]schema.Index{
					{Name: "index_foo", Columns: []string{"foo", "bar"}, Unique: true, Primary: true},
				}, nil).Once()
				mockSchema.EXPECT().GetForeignKeys("test").Return([]schema.ForeignKey{
					{
						Name:           "fk_foo",
						Columns:        []string{"foo"},
						ForeignTable:   "bar",
						ForeignColumns: []string{"baz"},
						OnDelete:       "cascade",
						OnUpdate:       "restrict",
					},
				}, nil).Once()
				for i := range successCaseExpected {
					mockContext.EXPECT().TwoColumnDetail(successCaseExpected[i][0], successCaseExpected[i][1]).Once()
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			beforeEach()
			test.setup()
			command := NewTableCommand(mockConfig, mockSchema)
			assert.NoError(t, command.Handle(mockContext))
		})
	}
}
