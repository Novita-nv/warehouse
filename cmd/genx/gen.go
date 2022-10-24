package genx

import (
	"fmt"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/pkg/structgen"
)

func Gen() {
	cfg := appctx.NewConfig()
	structgen.Create(structgen.Configuration{
		DbHost:     fmt.Sprintf("%s:%d", cfg.WriteDB.Host, cfg.WriteDB.Port),
		DbName:     cfg.WriteDB.Name,
		DbUser:     cfg.WriteDB.User,
		DbPassword: cfg.WriteDB.Pass,
		TagLabel:   "db,json",
	})
}

func GenLogic() {
	structgen.CreateLogic()
}
