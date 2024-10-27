function createLockerMap(lockers) {
   let grid = document.querySelector(".big-grid") 
   for(let locker of lockers) {
    let indicator = document.createElement("div")
    let shout = document.createElement("div")
    shout.classList.add("shout-message")
    shout.textContent = "Book available locker.";
    indicator.classList.add("indicator")
    let newGridItem = document.createElement("div") 
    newGridItem.classList.add("grid-item")
    let isOccupied = locker.userid !== ""
    if(isOccupied) {
        newGridItem.classList.add("occupied")
    }
    newGridItem.setAttribute("id", locker.lockernum)
    newGridItem.innerText = locker.lockernum
    newGridItem.appendChild(indicator)
    
    grid.appendChild(newGridItem)
    if(!isOccupied){
        newGridItem.addEventListener("mouseenter", (e) => {
            document.body.appendChild(shout); 
            shout.style.display = "block";
        });

        newGridItem.addEventListener("mousemove", (e) => {
            shout.style.left = `${e.pageX + 10}px`; 
            shout.style.top = `${e.pageY + 10}px`;  
        });

        newGridItem.addEventListener("mouseleave", () => {
            shout.style.display = "none";
            if (shout.parentElement) {
                shout.parentElement.removeChild(shout); 
            }
        });
    }
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