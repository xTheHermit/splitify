package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var _ ServerInterface = (*Server)(nil)

// (POST /api/auth/register)
func (h *Server) Register(c *gin.Context) {
	var input RegisterJSONRequestBody
	if err := c.Bind(&input); err != nil {
		c.String(http.StatusBadRequest, "Cannot parse JSON")
		return
	}

	success, err := h.authService.Register(input.Username, input.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while registering")
		return
	}

	if !success {
		c.String(http.StatusBadRequest, "Cannot register")
		return
	}

	c.String(http.StatusOK, "Success")
}

// (POST /api/auth/login)
func (h *Server) Login(c *gin.Context) {
	var input LoginJSONRequestBody
	if err := c.Bind(&input); err != nil {
		c.String(http.StatusBadRequest, "Cannot parse JSON")
		return
	}

	success, token, err := h.authService.Login(input.Username, input.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while logging in")
		return
	}

	if !success {
		c.String(http.StatusBadRequest, "Wrong username or password")
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"session_id",
		token,
		60*60*24*3,
		"/api",
		"",
		true,
		true,
	)
	c.Status(http.StatusOK)
}


// (GET /api/groups)
func (h *Server) GetMyGroups(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Not Implemented")
}

// (POST /api/groups)
func (h *Server) CreateGroup(c *gin.Context) {
	var input CreateGroupJSONRequestBody
	if err := c.Bind(&input); err != nil {
		c.String(http.StatusBadRequest, "Cannot parse JSON")
		return
	}

	id, err := h.groupService.Create(input.Name)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while creating group")
		return
	}

	c.JSON(http.StatusOK, Group{
		Id: id,
		Name: input.Name,
	})
}

// (DELETE /api/groups/{id})
func (h *Server) DeleteGroup(c *gin.Context, id int) {
	out, err := h.groupService.Delete(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Cannot delete group")
		return
	}

	if out {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

// (GET /api/groups/{id})
func (h *Server) GetGroup(c *gin.Context, id int) {
	out, err := h.groupService.Get(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Cannot get group")
		return
	}

	c.JSON(http.StatusOK, Group{
		Id: out.Id,
		Name: out.Name,
	})
}

// (PATCH /api/groups/{id})
func (h *Server) UpdateGroup(c *gin.Context, id int) {
	var input UpdateGroupJSONRequestBody
	if err := c.Bind(&input); err != nil {
		c.String(http.StatusBadRequest, "Cannot parse JSON")
		return
	}

	success, group, err := h.groupService.Update(id, input.Name)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while updating group")
		return
	}

	if success {
		c.JSON(http.StatusOK, Group{
			Id: group.Id,
			Name: group.Name,
		})
	} else {
		c.Status(http.StatusBadRequest)
	}
}

