window.onload = () => {
  drawChart();
  setInterval(drawChart, 1000 * 60 * 5);
};

const drawChart = async () => {
  const response = await fetch("/data/oneday");
  const datas = await response.json();
  let columns = [];
  let columns2 = [];
  let timeStamp = ["x"];
  let temperature = ["temperature"];
  let humidity = ["humidity"];
  let pressure = ["pressure"];
  for (let i = 0; i < datas.length; i++) {
    temperature.push(datas[i].temperature);
    humidity.push(datas[i].humidity);
    pressure.push(datas[i].pressure);

    const d = new Date(datas[i].created_at);
    // HACK:データベースにJSTをUTCとして入れてしまっている
    const formattedTime = `${d.getUTCFullYear()}-${d
      .getUTCMonth()
      .toString()
      .padStart(2, "0")}-${d.getUTCDate().toString().padStart(2, "0")} ${d
      .getUTCHours()
      .toString()
      .padStart(2, "0")}:${d.getUTCMinutes().toString().padStart(2, "0")}:${d
      .getUTCSeconds()
      .toString()
      .padStart(2, "0")}`;

    timeStamp.push(formattedTime);
  }
  columns.push(timeStamp, temperature, humidity);
  columns2.push(timeStamp, pressure);

  let chart = c3.generate({
    bindto: "#chart",
    data: {
      x: "x",
      xFormat: "%Y-%m-%d %H:%M:%S",
      columns: columns,
      axes: {
        humidity: "y2",
      },
    },
    axis: {
      x: {
        type: "timeseries",
        tick: {
          fit: false,
          format: "%H:%M",
        },
      },
      y2: {
        show: true,
      },
    },
    point: {
      show: false,
    },
  });

  let chart2 = c3.generate({
    bindto: "#chart2",
    data: {
      x: "x",
      xFormat: "%Y-%m-%d %H:%M:%S",
      columns: columns2,
    },
    axis: {
      x: {
        type: "timeseries",
        tick: {
          fit: false,
          format: "%H:%M",
        },
      },
    },
    point: {
      show: false,
    },
  });
};
