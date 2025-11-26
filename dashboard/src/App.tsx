import React, { ChangeEvent, useEffect, useState } from 'react';
import './App.css';
import { Advanced } from './Advanced';
import { Graph } from './Graph';
import { Logs } from './Logs';

type BrainState = {
  PollDelayMs: number;
  StateOfBoiler: string;
  CalculatedBoilerState: string;
  SmartSwitchOn: boolean;
  TemperatureCelsius: number;
  ThermostatThresholdCelsius: number;
  BoilerMode: string;
};

enum ViewMode {
  None = 0,
  Graph,
  Logs,
  Advanced,
}

const viewModes = [ViewMode.Graph, ViewMode.Logs, ViewMode.Advanced];

const getViewModeName = (mode: ViewMode): string => {
  switch (mode) {
    case ViewMode.Graph:
      return "Graph";
    case ViewMode.Logs:
      return "Logs";
    case ViewMode.Advanced:
      return "Advanced";
  }
  return "";
}

function App() {
  const [brainState, setBrainState] = useState<BrainState | null>(null);
  const [thermostat, setThermostat] = useState<number | "">("");
  const [boilerMode, setBoilerMode] = useState<string>("");
  const [thermostatUpdated, setThermostatUpdated] = useState<boolean>(false);
  const [viewMode, setViewMode] = useState<ViewMode>(ViewMode.None);

  const updateState = (setThermostatInput: boolean) => {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", "http://192.168.86.100:8080/update-thermostat/", false);
    xmlHttp.send(null);
    console.log(xmlHttp.responseText);
    var data: BrainState = JSON.parse(xmlHttp.responseText);
    setBrainState(data);
    if (setThermostatInput) {
      // TODO: Deal with what happens if another user changes the thermostat while the page is open
      setThermostat(data.ThermostatThresholdCelsius);
      setBoilerMode(data.BoilerMode)
    }
    setTimeout(() => updateState(false), data.PollDelayMs);
  }

  useEffect(() => {
    updateState(true);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);


  function saveThermostatThreshold(event: ChangeEvent<HTMLInputElement>): void {
    setThermostat(+event.currentTarget.value);
  }

  function setNewThermostatValue(): void {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", "http://192.168.86.100:8080/update-thermostat/?threshold=" + thermostat + "&mode=" + boilerMode, false);
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
          State Of Boiler: {brainState?.StateOfBoiler}
          {brainState?.StateOfBoiler !== brainState?.CalculatedBoilerState && (
            <>(Will change to {brainState?.CalculatedBoilerState} on next update)</>
          )}
        </div>
        <div>
          Smart Switch State: {brainState?.SmartSwitchOn ? "on" : "off"}
        </div>
        <div>
          Boiler Mode:
          <button disabled={boilerMode === "off"} onClick={() => setBoilerMode("off")}>Off</button>
          <button disabled={boilerMode === "auto"} onClick={() => setBoilerMode("auto")}>Auto</button>
          <button disabled={boilerMode === "on"} onClick={() => setBoilerMode("on")}>On</button>
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
            disabled={(brainState?.ThermostatThresholdCelsius || "") === thermostat && (brainState?.BoilerMode) === boilerMode}
            onClick={setNewThermostatValue}
          >
            Set
          </button>
          {thermostatUpdated && (<span>Saved!</span>)}
        </div>

        {viewModes.map(m => {
          if (viewMode === m) {
            return <button className='App-button_mode' onClick={() => setViewMode(ViewMode.None)}>Hide {getViewModeName(m)}</button>;
          } else {
            return <button className='App-button_mode' onClick={() => setViewMode(m)}>Show {getViewModeName(m)}</button>;
          }
        })}

        {viewMode !== ViewMode.None && (
          <div className='App-div_view-mode-box'>
            {viewMode === ViewMode.Graph && (
              <Graph />
            )}

            {viewMode === ViewMode.Logs && (
              <Logs />
            )}

            {viewMode === ViewMode.Advanced && (
              <Advanced />
            )}
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
