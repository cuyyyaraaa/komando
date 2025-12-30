package auth

import (
	"database/sql"

	"komando/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// Login godoc
// @Summary Login
// @Description Login by NIP and password, returns JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "login payload"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	res, err := h.svc.Login(req.NIP, req.Password)
	if err != nil {
		// if nip not found => sql.ErrNoRows
		if err == sql.ErrNoRows {
			response.Fail(c, 401, "UNAUTHORIZED", "invalid credentials")
			return
		}
		// bcrypt mismatch etc -> invalid credentials
		response.Fail(c, 401, "UNAUTHORIZED", "invalid credentials")
		return
	}

	response.OK(c, res)
}

// Me godoc
// @Summary Current user profile
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	uidAny, ok := c.Get("user_id")
	if !ok {
		response.Fail(c, 401, "UNAUTHORIZED", "missing user")
		return
	}
	uid := uidAny.(int64)

	u, err := h.svc.Me(uid)
	if err != nil {
		response.Fail(c, 500, "INTERNAL_ERROR", "failed to fetch user")
		return
	}
	response.OK(c, u)
}
