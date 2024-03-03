import React, { ChangeEvent, useEffect, useState } from 'react';
import './App.css';
import { Advanced } from './Advanced';
import { Graph } from './Graph';
import { Logs } from './Logs';

type BrainState = {
  PollDelayMs: number;
  BoilerState: string;
  SmartSwitchOn: boolean;
  TemperatureCelsius: number;
  ThermostatThresholdCelsius: number;
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
      setThermostat(data.ThermostatThresholdCelsius)
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
