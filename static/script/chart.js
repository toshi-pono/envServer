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

  let alp = 0.8;

  let tempC = datas[0].temperature;
  let humC = datas[0].humidity;
  let presC = datas[0].pressure;
  let beforeDate = new Date(datas[0].created_at);
  // ローパスフィルターを通して出力
  for (let i = 0; i < datas.length; i++) {
    let nowDate = new Date(datas[i].created_at);
    // 前回計測が差3分以内だったらローパスフィルターを通す
    if (nowDate.getTime() - beforeDate.getTime() < 1000 * 60 * 3) {
      tempC = tempC * alp + datas[i].temperature * (1 - alp);
      humC = humC * alp + datas[i].humidity * (1 - alp);
      presC = presC * alp + datas[i].pressure * (1 - alp);
    } else {
      tempC = datas[i].temperature;
      humC = datas[i].humidity;
      presC = datas[i].pressure;
    }

    const formattedTime = `${nowDate.getUTCFullYear()}-${nowDate
      .getUTCMonth()
      .toString()
      .padStart(2, "0")}-${nowDate
      .getUTCDate()
      .toString()
      .padStart(2, "0")} ${nowDate
      .getUTCHours()
      .toString()
      .padStart(2, "0")}:${nowDate
      .getUTCMinutes()
      .toString()
      .padStart(2, "0")}:${nowDate
      .getUTCSeconds()
      .toString()
      .padStart(2, "0")}`;

    temperature.push(tempC.toFixed(2));
    humidity.push(humC.toFixed(2));
    pressure.push(presC.toFixed(2));
    timeStamp.push(formattedTime);
    beforeDate = nowDate;
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

// 分の10の位まで同じかどうか判定する
const isTimeNear = (nowDate, beforeDate) => {
  let nowD = new Date(nowDate);
  let beforeD = new Date(beforeDate);
  return (
    nowD.getUTCFullYear() == beforeD.getUTCFullYear() &&
    nowD.getUTCMonth() == beforeD.getUTCMonth() &&
    nowD.getUTCDate() == beforeD.getUTCDate() &&
    nowD.getUTCHours() == beforeD.getUTCHours() &&
    nowD.getUTCMinutes().toString().padStart(2, "0")[0] ==
      beforeD.getUTCMinutes().toString().padStart(2, "0")[0]
  );
};
