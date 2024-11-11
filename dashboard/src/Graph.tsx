import Chart from 'chart.js/auto';
import { useEffect, useRef, useState } from 'react';
import zoomPlugin from 'chartjs-plugin-zoom';
import 'chartjs-adapter-date-fns';
import './Graph.css';

Chart.register(zoomPlugin);

type TimePoint = {
    Time: number
    Value: number
}

type GraphDataResponse = {
    Temperature: TimePoint[]
    Temperature1: TimePoint[]
    Temperature2: TimePoint[]
    Thermostat: TimePoint[]
    SmartSwitchState: TimePoint[]
    BoilerState: TimePoint[]
}

const temperatureLine = (label: string, data: TimePoint[], color: string, stepped: boolean = false) => {
    return {
        label: label,
        data: (data ?? []).map((p) => {
            return { x: p.Time, y: p.Value }
        }),
        borderColor: color,
        borderWidth: 1,
        stepped: stepped,
        fill: false,
        yAxisID: 'y1',
    }
}

const onOffSwitch = (label: string, data: TimePoint[], color: string) => {
    return {
        label: label,
        data: (data ?? []).map((p) => {
            return { x: p.Time, y: p.Value }
        }),
        borderColor: [color],
        backgroundColor: color,
        borderWidth: 0.1,
        pointRadius: 0,
        stepped: true,
        fill: true,
        yAxisID: 'y2'
    }
}

export function Graph() {
    const chartContainer = useRef<HTMLCanvasElement>(null);
    const graphContainer = useRef<HTMLDivElement>(null);
    const [days, setDays] = useState(7)

    useEffect(() => {
        if (chartContainer.current) {
            let xmlHttp = new XMLHttpRequest();
            xmlHttp.open("GET", "http://192.168.86.100:8080/graph-data/?days=" + days, false);
            xmlHttp.send(null);
            console.log(xmlHttp.responseText);
            var data: GraphDataResponse = JSON.parse(xmlHttp.responseText);

            let myChart = new Chart<"line", number[] | { x: number, y: number }[]>(chartContainer.current, {
                type: 'line',
                data: {
                    labels: data.Temperature.map((p => p.Time)),
                    datasets: [
                        temperatureLine('Dining Room', data.Temperature, 'rgba(75, 192, 192, 1)'),
                        temperatureLine('Bedroom', data.Temperature1, 'yellow'),
                        temperatureLine('Lounge', data.Temperature2, 'orange'),
                        temperatureLine('Thermostat', data.Thermostat, 'gray', true),
                        onOffSwitch('Boiler On', data.BoilerState, 'rgba(255, 50, 50, 1)'), // red
                        onOffSwitch('Smart Switch On', data.SmartSwitchState, 'rgba(0, 0, 255, 0.25)'), // blue
                    ],
                },
                options: {
                    maintainAspectRatio: false,
                    scales: {
                        x: {
                            type: "time",
                            time: {
                                parser: 'HH:mm:ss',
                                displayFormats: {
                                    hour: 'd/M HH:mm',
                                    day: 'd MMM',
                                },
                                tooltipFormat: 'd MMM yyyy - HH:mm:ss'
                            },
                        },
                        y1: {
                        },
                        y2: {
                            display: false
                        },
                    },
                    plugins: {
                        zoom: {
                            pan: {
                                enabled: true,
                                mode: 'x',
                            },
                            zoom: {
                                wheel: {
                                    enabled: true,
                                },
                                pinch: {
                                    enabled: true,
                                },
                                drag: {
                                    enabled: true,
                                    modifierKey: 'shift',
                                },
                                mode: 'x',
                            },
                        },
                    },
                    interaction: {
                        mode: 'x'
                    }
                },
            });
            return () => {
                myChart.destroy();
            };
        }
    }, [days]);

    const toggleFullScreen = () => {
        if (document.fullscreenElement) {
            document.exitFullscreen()
        } else {
            graphContainer.current && graphContainer.current.requestFullscreen()
        }
    }

    return (<div className='graph' ref={graphContainer}>
        {timePeriodButton("1 day", 1)}
        {" - "}
        {timePeriodButton("7 days", 7)}
        {" - "}
        {timePeriodButton("1 month", 31)}
        {" - "}
        {timePeriodButton("6 months", 6 * 31)}
        {" --- "}
        <button onClick={() => toggleFullScreen()}>Fullscreen</button>
        <canvas ref={chartContainer} style={{ maxWidth: '2000px', maxHeight: '500px' }}></canvas>
    </div>);

    function timePeriodButton(label: string, desiredDays: number) {
        return <button onClick={() => setDays(desiredDays)} disabled={days === desiredDays}>{label}</button>;
    }
}
