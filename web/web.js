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
    xmlHttp.open( "GET", theUrl, false ); 
    xmlHttp.send( null );
    return xmlHttp.responseText;
}

document.addEventListener("DOMContentLoaded",main)


var initial

function continuousUpdate() {
    setInterval( () => {
        let resp = httpGet("http://127.0.0.1:8080/lockerStatus/")
        if (resp !== initial) {
            document.querySelector(".big-grid").innerHTML = ''
            createLockerMap(JSON.parse(resp))
            foo()
            continuousUpdate()
        }
    }, 1000)
}


function main() {
    res = httpGet("http://127.0.0.1:8080/lockerStatus/")
    console.log(res)
    createLockerMap(JSON.parse(res))
    foo()
  };

function foo(){
  document.querySelectorAll(".grid-item").forEach(item => {
    item.addEventListener("click", function(event) {
        console.log("clicked on grid-item")
        if (item.classList.contains("occupied")) {
            alert("This locker is already booked!");
        } else {
            window.location.href = "booklocker.html"; 
        }
    });
});
}