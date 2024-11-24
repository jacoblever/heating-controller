package stores

type Config struct {
	CurrentTemperatureFilePath     string
	TemperatureLogFilePath         string
	TemperatureLog1FilePath        string
	TemperatureLog2FilePath        string
	OutsideTemperatureLogFilePath  string
	ThermostatThresholdLogFilePath string
	SmartSwitchLastAliveFilePath   string
	SmartSwitchStateLogFilePath    string
	BoilerStateLogFilePath         string
	BoilerLogFilePath              string
}

var DefaultConfig Config = Config{
	CurrentTemperatureFilePath:     "./current-temperature.txt",
	TemperatureLogFilePath:         "./temperature-log.txt",
	TemperatureLog1FilePath:        "./temperature-log-1.txt",
	TemperatureLog2FilePath:        "./temperature-log-2.txt",
	OutsideTemperatureLogFilePath:  "./outside-temperature-log.txt",
	ThermostatThresholdLogFilePath: "./thermostat-threshold-log.txt",
	SmartSwitchLastAliveFilePath:   "./smart-switch-last-alive.txt",
	SmartSwitchStateLogFilePath:    "./smart-switch-state-log.txt",
	BoilerStateLogFilePath:         "./boiler-state-log.txt",
	BoilerLogFilePath:              "./boiler-log.txt",
}

func (c Config) AllFilePaths() []string {
	return []string{
		c.CurrentTemperatureFilePath,
		c.TemperatureLogFilePath,
		c.TemperatureLog1FilePath,
		c.TemperatureLog2FilePath,
		c.OutsideTemperatureLogFilePath,
		c.ThermostatThresholdLogFilePath,
		c.SmartSwitchLastAliveFilePath,
		c.SmartSwitchStateLogFilePath,
		c.BoilerStateLogFilePath,
		c.BoilerLogFilePath,
	}
}
