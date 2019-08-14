package handlers

import (
	"stash.tutu.ru/opscore-workshop-admin/request-generator/resources"
)

type Handler struct {
	Res *resources.Resources
}

func New(res *resources.Resources) (r *Handler) {
	return &Handler{Res: res}
}