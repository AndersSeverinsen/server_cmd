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
document.addEventListener("DOMContentLoaded", function(event){
    res = httpGet("http://127.0.0.1:8080/lockerStatus/")
    console.log(res)
    createLockerMap(JSON.parse(res))
  });

  
  