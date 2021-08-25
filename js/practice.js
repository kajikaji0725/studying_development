// ==UserScript==
// @name         Gakujo-seiseki-highlighter
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  色分け
// @author       わ　た　し
// @match        https://gakujo.shizuoka.ac.jp/kyoumu/seisekiStudentSort.do;jsessionid=t-hpPk2yDntBU5zsIc-dTF1cLBp94aOK2EVA1cDF?sortColumn=shuutokuDate
// @icon         https://www.google.com/s2/favicons?domain=shizuoka.ac.jp
// @grant        none
// ==/UserScript==

(function() {
    'use strict';
    let rows = document.querySelectorAll("body > table:nth-child(10) > tbody > tr > td > table > tbody > tr")
    rows = Array.prototype.slice.call(rows,2)

    for(let row of rows){
        const grade = row[1].querySelector("td:nth-child(6)").textContent.trim();
        row.onmouseober = "";
        row.onmouseout = ""
        switch (grade){
            case "秀":
                row.style.backgroundColor = "red"
                break;
            case "優":
                row.style.backgroundColor = "bule"
                break;
            case "良":
                row.style.backgroundColor = "yellow"
                break;
            case "可":
                row.style.backgroundColor = "green"
                break;
            case "不可":
                row.style.backgroundColor = "black"
                break;
            default:
                console.error(`Unkown grade : ${grade}`);
        }
    }
})();