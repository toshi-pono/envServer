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

  let cnt = 0;
  let before = datas[0];
  let humiditySum = 0;
  let temperatureSum = 0;
  let pressureSum = 0;
  let minutesSum = 0;
  // 10分間の平均を出力
  for (let i = 0; i < datas.length; i++) {
    temperatureSum += datas[i].temperature;
    humiditySum += datas[i].humidity;
    pressureSum += datas[i].pressure;
    minutesSum += new Date(datas[i].created_at).getUTCMinutes();
    cnt++;
    if (
      !isTimeNear(datas[i].created_at, before.created_at) ||
      datas.length < 10
    ) {
      // 平均を追加
      temperature.push(temperatureSum / cnt);
      humidity.push(humiditySum / cnt);
      pressure.push(pressureSum / cnt);

      // 時刻の処理
      const d = new Date(datas[i].created_at);
      // HACK:データベースにJSTをUTCとして入れてしまっている
      const formattedTime = `${d.getUTCFullYear()}-${d
        .getUTCMonth()
        .toString()
        .padStart(2, "0")}-${d.getUTCDate().toString().padStart(2, "0")} ${d
        .getUTCHours()
        .toString()
        .padStart(2, "0")}:${Math.round(minutesSum / cnt)
        .toString()
        .padStart(2, "0")}:${d.getUTCSeconds().toString().padStart(2, "0")}`;
      timeStamp.push(formattedTime);

      // 変数初期化
      humiditySum = 0;
      temperatureSum = 0;
      pressureSum = 0;
      minutesSum = 0;
      cnt = 0;
    }

    before = datas[i];
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
