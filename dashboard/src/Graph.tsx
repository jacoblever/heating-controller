import Chart from 'chart.js/auto';
import { useEffect, useRef } from 'react';
import zoomPlugin from 'chartjs-plugin-zoom';
import 'chartjs-adapter-date-fns';

Chart.register(zoomPlugin);

type TimePoint = {
    Time: number
    Value: number
}

type GraphDatResponse = {
    Temperature: TimePoint[]
}

export function Graph() {
    const chartContainer = useRef<HTMLCanvasElement>(null);

    useEffect(() => {
        if (chartContainer.current) {
            let xmlHttp = new XMLHttpRequest();
            xmlHttp.open("GET", "http://192.168.86.100:8080/graph-data/", false);
            xmlHttp.send(null);
            console.log(xmlHttp.responseText);
            var data: GraphDatResponse = JSON.parse(xmlHttp.responseText);

            let myChart = new Chart(chartContainer.current, {
                type: 'line',
                data: {
                    labels: data.Temperature.map((p => p.Time)),
                    datasets: [
                        {
                            label: `Temperature`,
                            data: data.Temperature.map((p) => p.Value),
                            borderColor: 'rgba(75, 192, 192, 1)',
                            borderWidth: 1,
                            fill: false,
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

                        }
                    },
                    plugins: {
                        zoom: {
                            pan: {
                                enabled: true,
                                modifierKey: 'ctrl',
                            },
                            zoom: {
                                drag: {
                                    enabled: true
                                },
                                mode: 'x',
                            },
                        },
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
