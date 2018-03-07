package postgres

import (
	"context"

	"database/sql"

	"strings"

	rstypes "git.containerum.net/ch/json-types/resource-service"
	"git.containerum.net/ch/kube-client/pkg/cherry/resource-service"
	kubtypes "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func (db *pgDB) createServicePorts(ctx context.Context, serviceID, domain string,
	serviceType rstypes.ServiceType, ports []kubtypes.ServicePort) (err error) {
	db.log.WithField("service_id", serviceID).Debugf("add service ports %#v", ports)

	var query string
	switch serviceType {
	case rstypes.ServiceInternal:
		query = /* language=sql */ `INSERT INTO service_ports
			(service_id, name, port, target_port, protocol, domain_id)
			VALUES (:service_id, :name, :port, :target_port, :protocol, NULL)`
	case rstypes.ServiceExternal:
		query = /* language=sql */ `INSERT INTO service_ports
			(service_id, name, port, target_port, protocol, domain_id)
			SELECT :service_id, :name, :port, :target_port, :protocol, d.id
			FROM domains d
			WHERE d.domain = :domain`
	}

	stmt, err := db.preparer.PrepareNamed(query)
	if err != nil {
		err = rserrors.ErrDatabase().Log(err, db.log)
		return
	}
	defer stmt.Close()

	for _, port := range ports {
		_, err = stmt.ExecContext(ctx, rstypes.Port{
			ServiceID:  serviceID,
			Name:       port.Name,
			Port:       port.Port,
			TargetPort: port.TargetPort,
			Protocol:   rstypes.PortProtocol(strings.ToLower(string(port.Protocol))),
			Domain:     &domain,
		})
		if err != nil {
			err = rserrors.ErrDatabase().Log(err, db.log)
			return
		}
	}

	return
}

func (db *pgDB) CreateService(ctx context.Context, userID, nsLabel string, serviceType rstypes.ServiceType, req kubtypes.Service) (err error) {
	db.log.WithFields(logrus.Fields{
		"type":     serviceType,
		"user_id":  userID,
		"ns_label": nsLabel,
	}).Debugf("create service %#v", req)

	nsID, err := db.getNamespaceID(ctx, userID, nsLabel)
	if err != nil {
		return
	}
	if nsID == "" {
		err = rserrors.ErrResourceNotExists().Log(err, db.log)
		return
	}

	var deplID string
	query, args, _ := sqlx.Named( /* language=sql */
		`SELECT id FROM deployments WHERE ns_id = :ns_id AND name = :name`,
		map[string]interface{}{"ns_id": nsID, "name": req.Deploy})
	err = sqlx.GetContext(ctx, db.extLog, &deplID, db.extLog.Rebind(query), args...)
	switch err {
	case nil:
	case sql.ErrNoRows:
		err = rserrors.ErrResourceNotExists().Log(err, db.log)
		return
	default:
		err = rserrors.ErrDatabase().Log(err, db.log)
		return
	}

	var serviceID string
	query, args, _ = sqlx.Named( /* language=sql */
		`INSERT INTO services
		(deploy_id, name, type)
		VALUES (:deploy_id, :name, :type)
		RETURNING id`,
		map[string]interface{}{"deploy_id": deplID, "name": req.Name, "type": serviceType})
	err = sqlx.GetContext(ctx, db.extLog, &serviceID, db.extLog.Rebind(query), args...)
	if err != nil {
		err = rserrors.ErrDatabase().Log(err, db.log)
		return
	}

	_, err = sqlx.NamedExecContext(ctx, db.extLog, /* language=sql */
		`INSERT INTO permissions
		(kind, resource_id, resource_label, owner_user_id, user_id)
		VALUES (
			(CASE :service_type
				WHEN 'external' THEN CAST('extservice' AS RESOURCE_KIND)
				WHEN 'internal' THEN CAST('intservice' AS RESOURCE_KIND)
			END),
			:service_id,
			:service_name,
			:user_id,
			:user_id
		)`,
		map[string]interface{}{
			"service_type": serviceType,
			"service_id":   serviceID,
			"service_name": req.Name,
			"user_id":      userID,
		})
	if err != nil {
		err = rserrors.ErrDatabase().Log(err, db.log)
		return
	}

	err = db.createServicePorts(ctx, serviceID, req.Domain, serviceType, req.Ports)
	return
}

func (db *pgDB) getRawServices(ctx context.Context, nsID string) (serviceMap map[string]kubtypes.Service, serviceIDs []string, err error) {
	db.log.WithField("ns_id", nsID).Debug("get raw services")

	dbEntries := make([]rstypes.Service, 0)
	query, args, _ := sqlx.Named( /* language=sql */
		`SELECT 
			s.id,
			d.name AS depl_id,
			s.name,
			s.type,
			s.created_at,
			s.deleted,
			s.delete_time
		FROM services s
		JOIN deployments d ON s.deploy_id = d.id
		WHERE NOT s.deleted`,
		map[string]interface{}{"ns_id": nsID})
	err = sqlx.SelectContext(ctx, db.extLog, &dbEntries, db.extLog.Rebind(query), args...)
	if err != nil {
		err = rserrors.ErrDatabase().Log(err, db.log)
		return
	}

	serviceMap = make(map[string]kubtypes.Service)
	for _, v := range dbEntries {
		createdAt := v.CreatedAt.Unix()
		serviceMap[v.ID] = kubtypes.Service{
			Name:      v.Name,
			CreatedAt: &createdAt,
			Deploy:    v.DeployID,
		}
		serviceIDs = append(serviceIDs, v.ID)
	}

	return
}

func (db *pgDB) getServicesPorts(ctx context.Context, serviceIDs []string, serviceMap map[string]kubtypes.Service) (err error) {
	db.log.Debugf("get services ports %#v", serviceIDs)

	dbEntries := make([]rstypes.Port, 0)
	query, args, _ := sqlx.In( /* language=sql */
		`SELECT
			sp.id,
			sp.service_id,
			sp.name,
			sp.port,
			sp.target_port,
			sp.protocol,
			d.domain	
		FROM service_ports sp
		LEFT JOIN domains d ON sp.domain_id = d.id
		WHERE sp.service_id IN (?)`, serviceIDs)
	err = sqlx.SelectContext(ctx, db.extLog, &dbEntries, db.extLog.Rebind(query), args...)
	if err != nil {
		err = rserrors.ErrDatabase().Log(err, db.log)
		return
	}

	portMap := make(map[string][]kubtypes.ServicePort)
	for _, v := range dbEntries {
		ports := portMap[v.ServiceID]
		ports = append(ports, kubtypes.ServicePort{
			Name:       v.Name,
			Port:       v.Port,
			TargetPort: v.TargetPort,
			Protocol:   kubtypes.Protocol(strings.ToUpper(string(v.Protocol))),
		})
		portMap[v.ServiceID] = ports
	}

	for _, v := range dbEntries {
		service := serviceMap[v.ServiceID]
		service.Ports = portMap[v.ServiceID]
		serviceMap[v.ServiceID] = service
	}

	return
}

func (db *pgDB) GetServices(ctx context.Context, userID, nsLabel string) (ret []kubtypes.Service, err error) {
	db.log.WithFields(logrus.Fields{
		"user_id":  userID,
		"ns_label": nsLabel,
	}).Debug("get services")

	nsID, err := db.getNamespaceID(ctx, userID, nsLabel)
	if err != nil {
		return
	}
	if nsID == "" {
		err = rserrors.ErrResourceNotExists().Log(err, db.log)
		return
	}

	serviceMap, serviceIDs, err := db.getRawServices(ctx, nsID)
	if err != nil {
		return
	}

	if err = db.getServicesPorts(ctx, serviceIDs, serviceMap); err != nil {
		return
	}

	ret = make([]kubtypes.Service, 0)
	for _, v := range serviceMap {
		ret = append(ret, v)
	}

	return
}

func (db *pgDB) GetService(ctx context.Context, userID, nsLabel, serviceLabel string) (ret kubtypes.Service, err error) {
	db.log.WithFields(logrus.Fields{
		"user_id":       userID,
		"ns_label":      nsLabel,
		"service_label": serviceLabel,
	}).Debug("get service")

	nsID, err := db.getNamespaceID(ctx, userID, nsLabel)
	if err != nil {
		return
	}
	if nsID == "" {
		err = rserrors.ErrResourceNotExists().Log(err, db.log)
		return
	}

	var serviceEntry rstypes.Service
	query, args, _ := sqlx.Named( /* language=sql */
		`SELECT 
			s.id,
			d.name AS depl_id,
			s.name,
			s.type,
			s.created_at,
			s.deleted,
			s.delete_time
		FROM services s
		JOIN deployments d ON s.deploy_id = d.id
		WHERE s.name = :name AND NOT s.deleted`,
		map[string]interface{}{"ns_id": nsID, "name": serviceLabel})
	err = sqlx.GetContext(ctx, db.extLog, &serviceEntry, db.extLog.Rebind(query), args...)
	switch err {
	case nil:
	case sql.ErrNoRows:
		err = rserrors.ErrResourceNotExists().Log(err, db.log)
		return
	default:
		err = rserrors.ErrDatabase().Log(err, db.log)
		return
	}

	serviceIDs := []string{serviceEntry.ID}
	createdAt := serviceEntry.CreatedAt.Unix()
	serviceMap := map[string]kubtypes.Service{
		serviceEntry.ID: {
			Name:      serviceEntry.Name,
			CreatedAt: &createdAt,
			Deploy:    serviceEntry.DeployID,
		},
	}

	if err = db.getServicesPorts(ctx, serviceIDs, serviceMap); err != nil {
		return
	}

	ret = serviceMap[serviceEntry.ID]
	return
}

func (db *pgDB) UpdateService(ctx context.Context, userID, nsLabel, serviceLabel string, newServiceType rstypes.ServiceType, req kubtypes.Service) (err error) {
	db.log.WithFields(logrus.Fields{
		"user_id":          userID,
		"ns_label":         nsLabel,
		"service_label":    serviceLabel,
		"new_service_type": newServiceType,
	}).Debugf("update service to %#v", req)

	nsID, err := db.getNamespaceID(ctx, userID, nsLabel)
	if err != nil {
		return
	}
	if nsID == "" {
		err = rserrors.ErrResourceNotExists().Log(err, db.log)
		return
	}

	var serviceID string
	query, args, _ := sqlx.Named( /* language=sql */
		`WITH service_to_update AS (
			SELECT s.id
			FROM services s
			JOIN deployments d ON s.deploy_id = d.id
			WHERE d.ns_id = :ns_id AND s.name = :name
		)
		UPDATE services
		SET "type" = :new_type
		WHERE id = (SELECT id FROM service_to_update)
		RETURNING id`,
		map[string]interface{}{"ns_id": nsID, "name": serviceLabel, "new_type": newServiceType})
	err = sqlx.GetContext(ctx, db.extLog, &serviceID, db.extLog.Rebind(query), args...)
	switch err {
	case nil:
	case sql.ErrNoRows:
		err = rserrors.ErrResourceNotExists().Log(err, db.log)
		return
	default:
		err = rserrors.ErrDatabase().Log(err, db.log)
		return
	}

	_, err = sqlx.NamedExecContext(ctx, db.extLog, /* language=sql */
		`DELETE FROM service_ports WHERE service_id = :service_id`,
		map[string]interface{}{"servce_id": serviceID})
	if err != nil {
		err = rserrors.ErrDatabase().Log(err, db.log)
		return
	}

	err = db.createServicePorts(ctx, serviceID, req.Domain, newServiceType, req.Ports)
	return
}

func (db *pgDB) DeleteService(ctx context.Context, userID, nsLabel, serviceLabel string) (err error) {
	db.log.WithFields(logrus.Fields{
		"user_id":       userID,
		"ns_label":      nsLabel,
		"service_label": serviceLabel,
	}).Debug("delete service")

	nsID, err := db.getNamespaceID(ctx, userID, nsLabel)
	if err != nil {
		return
	}
	if nsID == "" {
		err = rserrors.ErrResourceNotExists().Log(err, db.log)
		return
	}

	result, err := sqlx.NamedExecContext(ctx, db.extLog, /* language=sql */
		`WITH service_to_update AS (
			SELECT s.id
			FROM services s
			JOIN deployments d ON s.deploy_id = d.id
			WHERE d.ns_id = :ns_id AND s.name = :name
		)
		UPDATE services
		SET deleted = TRUE, delete_time = now() AT TIME ZONE 'UTC'
		WHERE id = (SELECT id FROM service_to_update)`,
		map[string]interface{}{"ns_id": nsID, "name": serviceLabel})
	if count, _ := result.RowsAffected(); count <= 0 {
		err = rserrors.ErrResourceNotExists().Log(err, db.log)
	}

	return
}
