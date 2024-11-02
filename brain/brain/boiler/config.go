package boiler

type Config struct {
	CurrentTemperatureFilePath         string
	TemperatureLogFilePath             string
	TemperatureLog1FilePath            string
	TemperatureLog2FilePath            string
	CurrentThermostatThresholdFilePath string
	SmartSwitchLastAliveFilePath       string
	BoilerStateLogFilePath             string
	BoilerLogFilePath                  string
}

var DefaultConfig Config = Config{
	CurrentTemperatureFilePath:         "./current-temperature.txt",
	TemperatureLogFilePath:             "./temperature-log.txt",
	TemperatureLog1FilePath:            "./temperature-log-1.txt",
	TemperatureLog2FilePath:            "./temperature-log-2.txt",
	CurrentThermostatThresholdFilePath: "./current-thermostat-threshold.txt",
	SmartSwitchLastAliveFilePath:       "./smart-switch-last-alive.txt",
	BoilerStateLogFilePath:             "./boiler-state-log.txt",
	BoilerLogFilePath:                  "./boiler-log.txt",
}

func (c Config) AllFilePaths() []string {
	return []string{
		c.CurrentTemperatureFilePath,
		c.TemperatureLogFilePath,
		c.TemperatureLog1FilePath,
		c.TemperatureLog2FilePath,
		c.CurrentThermostatThresholdFilePath,
		c.SmartSwitchLastAliveFilePath,
		c.BoilerStateLogFilePath,
		c.BoilerLogFilePath,
	}
}
