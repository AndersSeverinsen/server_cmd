function createLockerMap(lockers) {
   let grid = document.querySelector(".big-grid") 
   for(let locker of lockers) {
    let indicator = document.createElement("div")
    indicator.classList.add("indicator")
    let newGridItem = document.createElement("div") 
    newGridItem.classList.add("grid-item")
    if(locker.userid !== "") {
        newGridItem.classList.add("occupied")
    }
    newGridItem.setAttribute("id", locker.lockernum)
    newGridItem.innerText = locker.lockernum
    newGridItem.appendChild(indicator)
    grid.appendChild(newGridItem)
   }
}

function httpGet(theUrl) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", theUrl, false ); // false for synchronous request
    xmlHttp.send( null );
    return xmlHttp.responseText;
}

var initial

function continuousUpdate() {
    setInterval( () => {
        console.log("dfsdfsd")
        let resp = httpGet("http://127.0.0.1:8080/lockerStatus/")
        if (resp !== initial) {
            document.querySelector(".big-grid").innerHTML = ''
            createLockerMap(JSON.parse(resp))
            continuousUpdate()
        }
    }, 1000)
}


function main() {
    res = httpGet("http://127.0.0.1:8080/lockerStatus/")
    console.log(res)
    createLockerMap(JSON.parse(res))
    initial = res 
    continuousUpdate()
}


document.addEventListener("DOMContentLoaded", main)
  
