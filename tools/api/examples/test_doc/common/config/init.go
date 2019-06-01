package config

func init() {
	CheckLogConfig()

	CheckConsulConfig()
	CheckAWSConfig()
	CheckEnvConfig()
	CheckServiceConfig()
	CheckXRayConfig()
	CheckSessionConfig()
	CheckHttpConfig()
	CheckRemoteApiConfig()
}
