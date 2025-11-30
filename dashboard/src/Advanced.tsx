import React, { useState } from 'react';

export function Advanced() {
    const [customNumber, setCustomNumber] = useState<number>(250);
    const [commandSent, setCommandSent] = useState<boolean>(false);

    function turnBoiler(commandToSend: string) {
        var xmlHttp = new XMLHttpRequest();
        xmlHttp.open("GET", "http://192.168.86.100:8080/turn-boiler/?command=" + commandToSend, false);
        xmlHttp.send(null);
        console.log(xmlHttp.responseText);

        setCommandSent(true);
        setTimeout(() => {
            setCommandSent(false);
        }, 2000);
    }

    return (<>
        <div style={{ marginBottom: "10px" }}>
            Turn Boiler (+ve is clockwise, and On direction):
        </div>
        <div style={{ marginBottom: "10px" }}>
            <button
                onClick={() => turnBoiler(`-${customNumber}`)}>
                Off (-ve) direction
            </button>
            <input
                type="number"
                value={customNumber}
                onChange={(e) => setCustomNumber(+e.currentTarget.value)}
            />
            <button
                onClick={() => turnBoiler(customNumber.toString())}>
                On (+ve) direction
            </button>
        </div>
        <div style={{ marginBottom: "10px" }}>
            <button onClick={() => turnBoiler('-250')}>Off 250</button>
            <button onClick={() => turnBoiler('-100')}>Off 100</button>
            <button onClick={() => turnBoiler('-50')}>Off 50</button>
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <button onClick={() => turnBoiler('50')}>On 50</button>
            <button onClick={() => turnBoiler('100')}>On 100</button>
            <button onClick={() => turnBoiler('250')}>On 250</button>
        </div>
        <div style={{ marginBottom: "10px" }}>
            <button onClick={() => turnBoiler('turn-anticlockwise')}>Simulate Off</button>
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <button onClick={() => turnBoiler('turn-clockwise')}>Simulate On</button>
        </div>
        <div>
            {commandSent && (<span>Command Sent!</span>)}
        </div>
    </>
    );
}
