package handlers

import (
	"net/http"

	m "git.containerum.net/ch/resource-service/pkg/router/middleware"
	"git.containerum.net/ch/resource-service/pkg/server"
	kubtypes "github.com/containerum/kube-client/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type ServiceHandlers struct {
	server.ServiceActions
	*m.TranslateValidate
}

// swagger:operation GET /namespaces/{namespace}/services Service GetServicesList
// Get services list.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserIDHeader'
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserNamespaceHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
// responses:
//  '200':
//    description: services list
//    schema:
//      $ref: '#/definitions/ServicesResponse'
//  default:
//    $ref: '#/responses/error'
func (h *ServiceHandlers) GetServicesListHandler(ctx *gin.Context) {
	resp, err := h.GetServicesList(ctx.Request.Context(), ctx.Param("namespace"))
	if err != nil {
		ctx.AbortWithStatusJSON(h.HandleError(err))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// swagger:operation GET /namespaces/{namespace}/services/{service} Service GetService
// Get services list.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserIDHeader'
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserNamespaceHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
//  - name: service
//    in: path
//    type: string
//    required: true
// responses:
//  '200':
//    description: service
//    schema:
//     $ref: '#/definitions/ResourceService'
//  default:
//    $ref: '#/responses/error'
func (h *ServiceHandlers) GetServiceHandler(ctx *gin.Context) {
	resp, err := h.GetService(ctx.Request.Context(), ctx.Param("namespace"), ctx.Param("service"))

	if err != nil {
		ctx.AbortWithStatusJSON(h.HandleError(err))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// swagger:operation POST /namespaces/{namespace}/services Service CreateService
// Create service.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserIDHeader'
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserNamespaceHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
//  - name: body
//    in: body
//    schema:
//     $ref: '#/definitions/Service'
// responses:
//  '201':
//    description: service created
//    schema:
//     $ref: '#/definitions/ResourceService'
//  default:
//    $ref: '#/responses/error'
func (h *ServiceHandlers) CreateServiceHandler(ctx *gin.Context) {
	var req kubtypes.Service
	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(h.BadRequest(ctx, err))
		return
	}

	createdService, err := h.CreateService(ctx.Request.Context(), ctx.Param("namespace"), req)
	if err != nil {
		ctx.AbortWithStatusJSON(h.HandleError(err))
		return
	}

	ctx.JSON(http.StatusCreated, createdService)
}

// swagger:operation POST /import/services Service ImportServices
// Import services.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserIDHeader'
//  - $ref: '#/parameters/UserRoleHeader'
//  - name: body
//    in: body
//    schema:
//     $ref: '#/definitions/ServicesList'
// responses:
//  '202':
//    description: services imported
//    schema:
//      $ref: '#/definitions/ImportResponse'
//  default:
//    $ref: '#/responses/error'
func (h *ServiceHandlers) ImportServicesHandler(ctx *gin.Context) {
	var req kubtypes.ServicesList
	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(h.BadRequest(ctx, err))
		return
	}

	resp := kubtypes.ImportResponse{
		Imported: []kubtypes.ImportResult{},
		Failed:   []kubtypes.ImportResult{},
	}

	for _, svc := range req.Services {
		if err := h.ImportService(ctx.Request.Context(), svc.Namespace, svc); err != nil {
			logrus.Warn(err)
			resp.ImportFailed(svc.Name, svc.Namespace, err.Error())
		} else {
			resp.ImportSuccessful(svc.Name, svc.Namespace)
		}
	}

	ctx.JSON(http.StatusAccepted, resp)
}

// swagger:operation PUT /namespaces/{namespace}/services/{service} Service UpdateService
// Update service.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserIDHeader'
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserNamespaceHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
//  - name: service
//    in: path
//    type: string
//    required: true
//  - name: body
//    in: body
//    schema:
//     $ref: '#/definitions/Service'
// responses:
//  '202':
//    description: service updated
//    schema:
//     $ref: '#/definitions/ResourceService'
//  default:
//    $ref: '#/responses/error'
func (h *ServiceHandlers) UpdateServiceHandler(ctx *gin.Context) {
	var req kubtypes.Service
	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(h.BadRequest(ctx, err))
		return
	}

	req.Name = ctx.Param("service")
	updatedService, err := h.UpdateService(ctx.Request.Context(), ctx.Param("namespace"), req)
	if err != nil {
		ctx.AbortWithStatusJSON(h.HandleError(err))
		return
	}

	ctx.JSON(http.StatusAccepted, updatedService)
}

// swagger:operation DELETE /namespaces/{namespace}/services/{service} Service DeleteService
// Delete service.
//
// ---
// x-method-visibility: public
// parameters:
//  - $ref: '#/parameters/UserIDHeader'
//  - $ref: '#/parameters/UserRoleHeader'
//  - $ref: '#/parameters/UserNamespaceHeader'
//  - name: namespace
//    in: path
//    type: string
//    required: true
//  - name: service
//    in: path
//    type: string
//    required: true
// responses:
//  '202':
//    description: service deleted
//  default:
//    $ref: '#/responses/error'
func (h *ServiceHandlers) DeleteServiceHandler(ctx *gin.Context) {
	err := h.DeleteService(ctx.Request.Context(), ctx.Param("namespace"), ctx.Param("service"))
	if err != nil {
		ctx.AbortWithStatusJSON(h.HandleError(err))
		return
	}

	ctx.Status(http.StatusAccepted)
}

// swagger:operation DELETE /namespaces/{namespace}/services Service DeleteAllServices
// Delete all services.
//
// ---
// x-method-visibility: private
// parameters:
//  - name: namespace
//    in: path
//    type: string
//    required: true
// responses:
//  '202':
//    description: all services in namespace deleted
//  default:
//    $ref: '#/responses/error'
func (h *ServiceHandlers) DeleteAllServicesHandler(ctx *gin.Context) {
	err := h.DeleteAllServices(ctx.Request.Context(), ctx.Param("namespace"))
	if err != nil {
		ctx.AbortWithStatusJSON(h.HandleError(err))
		return
	}

	ctx.Status(http.StatusAccepted)
}

// swagger:operation DELETE /namespaces/{namespace}/solutions/{solution}/services Service DeleteAllSolutionServices
// Delete all solution services.
//
// ---
// x-method-visibility: private
// parameters:
//  - name: namespace
//    in: path
//    type: string
//    required: true
//  - name: solution
//    in: path
//    type: string
//    required: true
// responses:
//  '202':
//    description: all solution services deleted
//  default:
//    $ref: '#/responses/error'
func (h *ServiceHandlers) DeleteAllSolutionServicesHandler(ctx *gin.Context) {
	if err := h.DeleteAllSolutionServices(ctx.Request.Context(), ctx.Param("namespace"), ctx.Param("solution")); err != nil {
		ctx.AbortWithStatusJSON(h.HandleError(err))
		return
	}

	ctx.Status(http.StatusAccepted)
}
