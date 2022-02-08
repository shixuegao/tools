// var global_ip = "10.116.192.91"
var global_ip = "localhost"
var global_port = 65001

function globalUri() {
  return "http://" + global_ip + ":" + global_port
}

//token
var tokenSignature = "pressureToolToken"
function setToken(token) {
  localStorage.setItem(tokenSignature, token)
}
function getToken() {
  return localStorage.getItem(tokenSignature)
}

//日期转换
function dateFormat(fmt, date) {
  let ret;
  const opt = {
      "Y+": date.getFullYear().toString(),        // 年
      "m+": (date.getMonth() + 1).toString(),     // 月
      "d+": date.getDate().toString(),            // 日
      "H+": date.getHours().toString(),           // 时
      "M+": date.getMinutes().toString(),         // 分
      "S+": date.getSeconds().toString()          // 秒
      // 有其他格式化字符需求可以继续添加，必须转化成字符串
  };
  for (let k in opt) {
      ret = new RegExp("(" + k + ")").exec(fmt);
      if (ret) {
          fmt = fmt.replace(ret[1], (ret[1].length == 1) ? (opt[k]) : (opt[k].padStart(ret[1].length, "0")))
      };
  };
  return fmt;
}

//IP检测
var regex = /^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$/
function legalIp(ip) {
  if (ip == "localhost" || regex.test(ip)) {
    return true
  }
  return false
}