package impl

import (
	"context"

	"git.containerum.net/ch/auth/proto"
	rstypes "git.containerum.net/ch/json-types/resource-service"
	"git.containerum.net/ch/kube-client/pkg/cherry/resource-service"
	"git.containerum.net/ch/resource-service/pkg/models"
	"git.containerum.net/ch/resource-service/pkg/server"
	"git.containerum.net/ch/utils"
	"github.com/sirupsen/logrus"
)

func (rs *resourceServiceImpl) SetUserAccesses(ctx context.Context, accessLevel rstypes.PermissionStatus) error {
	userID := utils.MustGetUserID(ctx)
	rs.log.WithFields(logrus.Fields{
		"user_id":      userID,
		"access_level": accessLevel,
	}).Info("set user resources access level")

	err := rs.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		if updErr := tx.SetAllResourcesAccess(ctx, userID, accessLevel); updErr != nil {
			return updErr
		}

		if updErr := rs.updateAccess(ctx, tx, userID); updErr != nil {
			return updErr
		}

		return nil
	})

	return err
}

func (rs *resourceServiceImpl) SetUserVolumeAccess(ctx context.Context, label string, req *rstypes.SetNamespaceAccessRequest) error {
	userID := utils.MustGetUserID(ctx)
	rs.log.WithFields(logrus.Fields{
		"user_id":      userID,
		"to":           req.Username,
		"label":        label,
		"access_level": req.Access,
	}).Info("change user volume access level")

	isAdmin := server.IsAdminRole(ctx)

	err := rs.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		vol, getErr := tx.GetUserVolumeByLabel(ctx, userID, label)
		if getErr != nil {
			return getErr
		}

		if vol.OwnerUserID != userID && !isAdmin {
			return rserrors.ErrResourceNotOwned()
		}

		if vol.Limited && !isAdmin {
			return rserrors.ErrPermissionDenied().AddDetailF("limited owner can`t assign permissions")
		}

		info, getErr := rs.User.UserInfoByLogin(ctx, req.Username)
		if getErr != nil {
			return getErr
		}

		vol.PermissionRecord.UserID = info.ID
		vol.PermissionRecord.AccessLevel = req.Access

		if setErr := tx.SetResourceAccess(ctx, &vol.PermissionRecord); setErr != nil {
			return setErr
		}

		if updErr := rs.updateAccess(ctx, tx, userID); updErr != nil {
			return updErr
		}

		return nil
	})

	return err
}

func (rs *resourceServiceImpl) SetUserNamespaceAccess(ctx context.Context, label string, req *rstypes.SetNamespaceAccessRequest) error {
	userID := utils.MustGetUserID(ctx)
	rs.log.WithFields(logrus.Fields{
		"user_id":      userID,
		"to":           req.Username,
		"label":        label,
		"access_level": req.Access,
	}).Info("change user volume access level")

	isAdmin := server.IsAdminRole(ctx)

	err := rs.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		ns, getErr := tx.GetUserNamespaceByLabel(ctx, userID, label)
		if getErr != nil {
			return getErr
		}

		if ns.OwnerUserID != userID && !isAdmin {
			return rserrors.ErrResourceNotOwned()
		}

		if ns.Limited && !isAdmin {
			return rserrors.ErrPermissionDenied().AddDetailF("limited owner can`t assign permissions")
		}

		info, getErr := rs.User.UserInfoByLogin(ctx, req.Username)
		if getErr != nil {
			return getErr
		}

		ns.PermissionRecord.UserID = info.ID
		ns.PermissionRecord.AccessLevel = req.Access

		if setErr := tx.SetResourceAccess(ctx, &ns.PermissionRecord); setErr != nil {
			return setErr
		}

		if updErr := rs.updateAccess(ctx, tx, userID); updErr != nil {
			return updErr
		}

		return nil
	})

	return err
}

func (rs *resourceServiceImpl) GetUserNamespaceAccesses(ctx context.Context, label string) (rstypes.GetUserNamespaceAccessesResponse, error) {
	userID := utils.MustGetUserID(ctx)
	rs.log.WithFields(logrus.Fields{
		"user_id": userID,
		"label":   label,
	}).Info("get user namespace accesses")

	ret, err := rs.DB.GetNamespaceWithUserPermissions(ctx, userID, label)

	return ret, err
}

func (rs *resourceServiceImpl) GetUserVolumeAccesses(ctx context.Context, label string) (rstypes.VolumeWithUserPermissions, error) {
	userID := utils.MustGetUserID(ctx)
	rs.log.WithFields(logrus.Fields{
		"user_id": userID,
		"label":   label,
	}).Info("get user volume accesses")

	ret, err := rs.DB.GetVolumeWithUserPermissions(ctx, userID, label)

	return ret, err
}

func (rs *resourceServiceImpl) GetUserAccesses(ctx context.Context) (*authProto.ResourcesAccess, error) {
	userID := utils.MustGetUserID(ctx)
	rs.log.WithField("user_id", userID).Info("get all user accesses")

	ret, err := rs.DB.GetUserResourceAccesses(ctx, userID)

	return ret, err
}

func (rs *resourceServiceImpl) DeleteUserNamespaceAccess(ctx context.Context, nsLabel string, req rstypes.DeleteNamespaceAccessRequest) error {
	rs.log.WithFields(logrus.Fields{
		"ns_label":    nsLabel,
		"target_user": req.Username,
	}).Info("delete user namespace access")

	err := rs.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		user, getErr := rs.User.UserInfoByLogin(ctx, req.Username)
		if getErr != nil {
			return getErr
		}

		ns, getErr := tx.GetUserNamespaceByLabel(ctx, user.ID, nsLabel)
		if getErr != nil {
			return getErr
		}

		if ns.OwnerUserID == ns.UserID {
			return rserrors.ErrDeleteOwnerAccess()
		}

		if delErr := tx.DeleteResourceAccess(ctx, ns.Resource, user.ID); delErr != nil {
			return delErr
		}

		if updErr := rs.updateAccess(ctx, tx, user.ID); updErr != nil {
			return updErr
		}

		return nil
	})

	return err
}

func (rs *resourceServiceImpl) DeleteUserVolumeAccess(ctx context.Context, volLabel string, req rstypes.DeleteVolumeAccessRequest) error {
	rs.log.WithFields(logrus.Fields{
		"vol_label":   volLabel,
		"target_user": req.Username,
	}).Info("delete user volume access")

	err := rs.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		user, getErr := rs.User.UserInfoByLogin(ctx, req.Username)
		if getErr != nil {
			return getErr
		}

		vol, getErr := tx.GetUserVolumeByLabel(ctx, user.ID, volLabel)
		if getErr != nil {
			return getErr
		}

		if vol.OwnerUserID == vol.UserID {
			return rserrors.ErrDeleteOwnerAccess()
		}

		if delErr := tx.DeleteResourceAccess(ctx, vol.Resource, user.ID); delErr != nil {
			return delErr
		}

		if updErr := rs.updateAccess(ctx, tx, user.ID); updErr != nil {
			return updErr
		}

		return nil
	})

	return err
}
