package jobs

//func (r *Resources) SyncLdapUsers(c context.Context, p hjob.Payload) error {
//	var in model.JobSyncLdapUsers
//	if err := p.Decode(&in); err != nil {
//		return tracer.Trace(err)
//	}
//
//	updated, err := r.app.SyncLdapUsers(c, in.Since)
//	if err != nil {
//		hlog.CtxLogger(c).Error("failed ldap sync", hlog.Err(tracer.Trace(err)))
//		return tracer.Trace(err)
//	}
//
//	hlog.CtxLogger(c).Info("synced ldap users", hlog.Int("updated", updated))
//	return nil
//}
