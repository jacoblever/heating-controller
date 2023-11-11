import React, { ChangeEvent, useEffect, useState } from 'react';
import './App.css';
import { Advanced } from './Advanced';

type BrainState = {
  PollDelayMs: number;
  BoilerState: string;
  SmartSwitchOn: boolean;
  TemperatureCelsius: number;
  ThermostatThresholdCelsius: number;
};

function App() {
  const [brainState, setBrainState] = useState<BrainState | null>(null);
  const [thermostat, setThermostat] = useState<number | "">("");
  const [thermostatUpdated, setThermostatUpdated] = useState<boolean>(false);
  const [showAdvanced, setShowAdvanced] = useState<boolean>(false);

  const updateState = (setThermostatInput: boolean) => {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", "http://192.168.86.100:8080/update-thermostat/", false);
    xmlHttp.send(null);
    console.log(xmlHttp.responseText);
    var data: BrainState = JSON.parse(xmlHttp.responseText);
    setBrainState(data);
    if (setThermostatInput) {
      // TODO: Deal with what happens if another user changes the thermostat while the page is open
      setThermostat(data.ThermostatThresholdCelsius)
    }
    setTimeout(() => updateState(false), data.PollDelayMs);
  }

  useEffect(() => {
    updateState(true);
  }, []);


  function saveThermostatThreshold(event: ChangeEvent<HTMLInputElement>): void {
    setThermostat(+event.currentTarget.value);
  }

  function setNewThermostatValue(): void {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", "http://192.168.86.100:8080/update-thermostat/?threshold=" + thermostat, false);
    xmlHttp.send(null);
    console.log(xmlHttp.responseText);
    var data = JSON.parse(xmlHttp.responseText);
    setBrainState(data);
    setThermostatUpdated(true);
    setTimeout(() => {
      setThermostatUpdated(false);
    }, 2000);
  }

  return (
    <div className="App">
      <header className="App-header">
        <h1>Heating Controller Dashboard</h1>
      </header>
      <div>
        <div>
          Boiler State: {brainState?.BoilerState}
        </div>
        <div>
          Smart Switch State: {brainState?.SmartSwitchOn ? "on" : "off"}
        </div>
        <div>
          Current Temperature: {brainState?.TemperatureCelsius}
        </div>
        <div>
          <label>
            Thermostat Threshold (Celsius):
            <input
              type="number"
              value={thermostat}
              onChange={saveThermostatThreshold}
            />
            <span id="thermostat-threshold-celsius-saved"></span>
          </label>
          <button
            disabled={(brainState?.ThermostatThresholdCelsius || "") === thermostat}
            onClick={setNewThermostatValue}
          >
            Set
          </button>
          {thermostatUpdated && (<span>Saved!</span>)}
        </div>

        {!showAdvanced && <button onClick={() => setShowAdvanced(true)}>Show Advanced</button>}
        {showAdvanced && <button onClick={() => setShowAdvanced(false)}>Close Advanced</button>}
        {showAdvanced && <Advanced />}
      </div>
    </div>
  );
}

export default App;
