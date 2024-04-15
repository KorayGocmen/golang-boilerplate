package context

func RemoteIP(ctx Ctx) string {
	remoteIP := ctx.Value(KeyRemoteIP)
	if remoteIP == nil {
		return ""
	}

	return remoteIP.(string)
}
