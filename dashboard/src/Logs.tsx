import { useEffect, useState } from 'react';

type LogLine = {
    Time: number
    Message: string
}

type LogsResponse = {
    Boiler: LogLine[]
    Brain: LogLine[]
}

export function Logs() {
    const [boilerLogs, setBoilerLogs] = useState<string>("");
    const [brainLogs, setBrainLogs] = useState<string>("");

    useEffect(() => {
        let xmlHttp = new XMLHttpRequest();
        xmlHttp.open("GET", "http://192.168.86.100:8080/logs/", false);
        xmlHttp.send(null);
        console.log(xmlHttp.responseText);
        var data: LogsResponse = JSON.parse(xmlHttp.responseText);
        setBoilerLogs(data.Boiler.map(x => x.Message).join("\n"));
        setBrainLogs(data.Brain.map(x => x.Message).join("\n"));
    }, []);

    return (<div>
        <h1>Boiler</h1>
        <textarea value={boilerLogs} rows={20} style={{ width: "80%" }}></textarea>
        <h1>Brain</h1>
        <textarea value={brainLogs} rows={20} style={{ width: "80%" }}></textarea>
    </div >);
}
