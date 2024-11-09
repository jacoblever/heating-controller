import Chart from 'chart.js/auto';
import { useEffect, useRef } from 'react';
import zoomPlugin from 'chartjs-plugin-zoom';
import 'chartjs-adapter-date-fns';

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

export function Graph() {
    const chartContainer = useRef<HTMLCanvasElement>(null);

    useEffect(() => {
        if (chartContainer.current) {
            let xmlHttp = new XMLHttpRequest();
            xmlHttp.open("GET", "http://192.168.86.100:8080/graph-data/", false);
            xmlHttp.send(null);
            console.log(xmlHttp.responseText);
            var data: GraphDataResponse = JSON.parse(xmlHttp.responseText);

            let myChart = new Chart<"line", number[] | { x: number, y: number }[]>(chartContainer.current, {
                type: 'line',
                data: {
                    labels: data.Temperature.map((p => p.Time)),
                    datasets: [
                        {
                            label: `Dining Room`,
                            data: (data.Temperature ?? []).map((p) => {
                                return { x: p.Time, y: p.Value }
                            }),
                            borderColor: 'rgba(75, 192, 192, 1)',
                            borderWidth: 1,
                            fill: false,
                            yAxisID: 'y1'
                        },
                        {
                            label: `Bedroom`,
                            data: (data.Temperature1 ?? []).map((p) => {
                                return { x: p.Time, y: p.Value }
                            }),
                            borderColor: 'yellow',
                            borderWidth: 1,
                            fill: false,
                            yAxisID: 'y1'
                        },
                        {
                            label: `Lounge`,
                            data: (data.Temperature2 ?? []).map((p) => {
                                return { x: p.Time, y: p.Value }
                            }),
                            borderColor: 'orange',
                            borderWidth: 1,
                            fill: false,
                            yAxisID: 'y1'
                        },
                        {
                            label: `Thermostat`,
                            data: (data.Thermostat ?? []).map((p) => {
                                return { x: p.Time, y: p.Value }
                            }),
                            borderColor: 'gray',
                            borderWidth: 1,
                            stepped: true,
                            fill: false,
                            yAxisID: 'y1'
                        },
                        {
                            label: 'Smart Switch State',
                            data: (data.SmartSwitchState ?? []).map((p) => {
                                return { x: p.Time, y: p.Value }
                            }),
                            borderColor: ['rgba(255, 99, 132, 1)'], //red
                            borderWidth: 0.1,
                            pointRadius: 0,
                            stepped: true,
                            fill: true,
                            yAxisID: 'y2'
                        },
                        {
                            label: 'Boiler State',
                            data: (data.BoilerState ?? []).map((p) => {
                                return { x: p.Time, y: p.Value === 1 ? 0.5 : 0 }
                            }),
                            borderColor: ['rgba(255, 99, 132, 1)'], //red
                            borderWidth: 0.1,
                            pointRadius: 0,
                            stepped: true,
                            fill: true,
                            yAxisID: 'y2'
                        },
                    ],
                },
                options: {
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
    }, []);

    return (<div>
        <canvas ref={chartContainer} style={{ maxWidth: '2000px', maxHeight: '500px' }}></canvas>
    </div>);
}
