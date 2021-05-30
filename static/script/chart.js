window.onload = () => {
  drawChart();
  setInterval(drawChart, 1000 * 60 * 5);
};

const drawChart = async () => {
  const response = await fetch("/data/oneday");
  const datas = await response.json();
  let columns = [];
  let timeStamp = ["x"];
  let temperature = ["temperature"];
  let humidity = ["humidity"];
  for (let i = 0; i < datas.length; i++) {
    temperature.push(datas[i].temperature);
    humidity.push(datas[i].humidity);

    const d = new Date(datas[i].created_at);
    const formattedTime = `${d.getFullYear()}-${d
      .getMonth()
      .toString()
      .padStart(2, "0")}-${d.getDate().toString().padStart(2, "0")} ${d
      .getHours()
      .toString()
      .padStart(2, "0")}:${d.getMinutes().toString().padStart(2, "0")}:${d
      .getSeconds()
      .toString()
      .padStart(2, "0")}`;

    timeStamp.push(formattedTime);
  }
  columns.push(timeStamp, temperature, humidity);

  let chart = c3.generate({
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
          format: "%Y-%m-%d %H:%M:%S",
          rotate: -90,
        },
      },
      y2: {
        show: true,
      },
    },
  });
};
