function Create() {
    let longURL = document.getElementById("longURL").value
    let shortURL = document.getElementById("shortURL").value
    let sendData
    let responseData
    if (shortURL == ""){
        sendData = JSON.stringify({
            longURL: longURL,
        });
    } else{
        sendData = JSON.stringify({
            longURL: longURL,
            shortURL: shortURL,
        });
    }
    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/create",true);
    xhr.send(sendData)
    xhr.onreadystatechange = function() {
        if (xhr.readyState != 4) return;
        if(xhr.status == 200) {
            document.getElementById("info").style.display = "block";
            responseData = JSON.parse(xhr.responseText);
            if (responseData.errorMsg == ""){
                document.getElementById("info").style.background = "palegreen";
                document.getElementById("info").innerHTML = ("longURL - " + responseData.longUrl +"<br>"+ "shortURL - " + responseData.shortUrl)
            } else{
                document.getElementById("info").style.background = "lightcoral";
                document.getElementById("info").innerHTML = (responseData.errorMsg)
            }

        }
        else {
            document.getElementById("info").style.background = "lightcoral";
            document.getElementById("info").innerHTML = ("error")
        }
    };
}
function Appear() {
    document.getElementById("customRoute").style.display = "block";
    document.getElementById("buttonCustomRoute").style.display = "none";
}