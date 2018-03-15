// Code generated by noice. DO NOT EDIT.
package rserrors

import (
	bytes "bytes"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry"
	template "text/template"
)

const ()

func ErrDatabase(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Database error", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x3, Kind: 0x1}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrResourceNotExists(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Resource is not exists", StatusHTTP: 404, ID: cherry.ErrID{SID: 0x3, Kind: 0x2}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrResourceAlreadyExists(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Resource already exists", StatusHTTP: 409, ID: cherry.ErrID{SID: 0x3, Kind: 0x3}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrPermissionDenied(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Permission denied", StatusHTTP: 403, ID: cherry.ErrID{SID: 0x3, Kind: 0x4}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrTariffUnchanged(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Tariff unchanged", StatusHTTP: 400, ID: cherry.ErrID{SID: 0x3, Kind: 0x5}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrTariffNotFound(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Tariff was not found", StatusHTTP: 404, ID: cherry.ErrID{SID: 0x3, Kind: 0x6}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrResourceNotOwned(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Can`t set access for resource which not owned by user", StatusHTTP: 403, ID: cherry.ErrID{SID: 0x3, Kind: 0x7}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrDeleteOwnerAccess(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Owner can`t delete has own access to resource", StatusHTTP: 409, ID: cherry.ErrID{SID: 0x3, Kind: 0x8}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrAccessRecordNotExists(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Access record for user not exists", StatusHTTP: 404, ID: cherry.ErrID{SID: 0x3, Kind: 0x9}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrInternal(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Internal error", StatusHTTP: 500, ID: cherry.ErrID{SID: 0x3, Kind: 0xa}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrValidation(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Validation error", StatusHTTP: 400, ID: cherry.ErrID{SID: 0x3, Kind: 0xb}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrServiceNotExternal(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Service is not external", StatusHTTP: 400, ID: cherry.ErrID{SID: 0x3, Kind: 0xc}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrTCPPortNotFound(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "TCP Port was not found in service", StatusHTTP: 400, ID: cherry.ErrID{SID: 0x3, Kind: 0xd}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrInsufficientStorage(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Can`t allocate resources for volume", StatusHTTP: 507, ID: cherry.ErrID{SID: 0x3, Kind: 0xe}, Details: []string(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}
func renderTemplate(templText string) string {
	buf := &bytes.Buffer{}
	templ, err := template.New("").Parse(templText)
	if err != nil {
		return err.Error()
	}
	err = templ.Execute(buf, map[string]string{})
	if err != nil {
		return err.Error()
	}
	return buf.String()
}
