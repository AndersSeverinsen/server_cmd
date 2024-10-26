



function updateLockermap(lockerJson) {
    lockers = lockerJson.lockers;
}



async function getLockerStatus() {
    const response = await fetch("/status/", {
        method: "GET", 
        headers: {Accept: "application/json"},
    })
    if (!response.ok) {
        throw new Error(
            `Connection failed with status code ${response.status}`
        );
    }
    log(response.json)
    return response.json;
}







