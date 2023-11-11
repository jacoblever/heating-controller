import React, { useState } from 'react';

export function Advanced() {
    const [command, setCommand] = useState<string>("");
    const [commandSent, setCommandSent] = useState<boolean>(false);

    function turnBoiler(commandOverride?: string) {
        var commandToSend = "";
        if (commandOverride) {
            commandToSend = commandOverride;
        } else {
            commandToSend = command;
            if (commandToSend === "") {
                return;
            }
        }

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
        <div>
            <label>
                Turn Boiler (+ve is clockwise):
                <input
                    type="text"
                    value={command}
                    onChange={(e) => setCommand(e.currentTarget.value)}
                />
                <button
                    onClick={() => turnBoiler()}>
                    Send
                </button>
                {commandSent && (<span>Sent!</span>)}
            </label>
        </div>
        <div>
            <button onClick={() => turnBoiler('turn-clockwise')}>Turn Clockwise</button>
            <button onClick={() => turnBoiler('turn-anticlockwise')}>Turn Anticlockwise</button>
        </div>
    </>
    );
}
