


function createLockers() {
    let grid = document.querySelector(".big-grid")
    for(let i = 0; i < 10; i++) {
        let indicator = document.createElement("div")
        indicator.classList.add("indicator")
        let newLocker = document.createElement("div")
        newLocker.classList.add("grid-item")
        newLocker.classList.add("occupied")
        newLocker.setAttribute("id", i)
        newLocker.innerText = i
        newLocker.appendChild(indicator)
        grid.appendChild(newLocker)
    }
}


function updateLockermap(lockerJson) {
    lockers = lockerJson.lockers;
}



function httpGet() {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", "http://localhost:8080/statusLocker", false ); // false for synchronous request
    xmlHttp.send( null );
    return xmlHttp.responseText;
}




document.addEventListener("DOMContentLoaded", function(event){
    createLockers()
    res = httpGet()
    console.log(res)
  });



