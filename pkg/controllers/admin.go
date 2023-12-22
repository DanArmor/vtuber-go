package controllers

import (
	"context"
	"net/http"

	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/pkg/controllers/resp"
	"github.com/DanArmor/vtuber-go/pkg/controllers/types"
	"github.com/gin-gonic/gin"
)

func extractOrgsNames(vtubers []types.Vtuber) []string {
	names := make([]string, 0, len(vtubers))
	for _, v := range vtubers {
		names = append(names, v.CompanyName)
	}
	return names
}
func (s *Service) orgsBuilders(names []string) []*ent.OrgCreate {
	builders := make([]*ent.OrgCreate, 0, len(names))
	for _, v := range names {
		builders = append(builders, s.Db.Org.Create().SetName(v))
	}
	return builders
}

func (s *Service) updateCompanies(vtubers []types.Vtuber) {
	names := extractOrgsNames(vtubers)
	err := s.Db.Org.
		CreateBulk(s.orgsBuilders(names)...).
		OnConflict().
		UpdateNewValues().
		Exec(context.Background())
}

func (s *Service) PostVtubers(c *gin.Context) {
	type PostVtubersInput struct {
		Vtubers []types.Vtuber
	}
	var input PostVtubersInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeCantBindJsonBody, "Can't bind json body"))
		return
	}

	c.JSON(http.StatusOK, resp.HandlerResult(gin.H{"vtubers": vtubers}))
}
