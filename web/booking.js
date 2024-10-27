document.querySelector("#button").addEventListener("click", function(event) {
    event.preventDefault();
    id = document.getElementById("name").value
    
    setTimeout(function() {
    window.location.href = "bookingconfirmed.html"; 
    },2000)

});

