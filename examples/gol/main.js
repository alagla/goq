
// function refresh(fun, millis){
//     fun();
//     setInterval(fun, millis);
// }



var hostEndpoint = "ws://localhost:8000/ws";
var socket = null;
var canvas = null;
var canvasOffsetLeft = null;
var canvasOffsetTop = null;
var golWidth = 81;
var golHeight = 81;
var cellStyle = "#000000";
var backgroundStyle = "#D3D3D3";
var cellW = 0;
var cellH = 0;

function main(){
    try{
        connectToHost();
    } catch (e) {
        console.error(e);
    }
    canvas = document.getElementById("GOLCanvas");
    let ctx = canvas.getContext("2d");

    cellW = canvas.width / golWidth;
    cellH = canvas.height / golHeight;
    canvasOffsetLeft = canvas.offsetLeft;
    canvasOffsetTop = canvas.offsetTop;

    clearCanvas(ctx);

    canvas.addEventListener('click', clickEvent);

    let nextGenButton = document.getElementById("nextGenButton");
    nextGenButton.addEventListener('click', nextGenEvent);
}

function clearCanvas(ctx) {
    ctx.fillStyle = backgroundStyle;
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    ctx.fillStyle = cellStyle;
}

function drawCell(ctx, x, y){
    ctx.fillRect(x * cellH, y * cellW, cellH, cellW);
}

function drawMap(themap){
    let ctx = canvas.getContext("2d");
    clearCanvas(ctx);
    for (let idx in themap){
        if (themap[idx] == 1){
            drawCell(ctx, idx % golWidth, idx / golWidth);
        }
    }
}

function connectToHost(){
    socket = new WebSocket(hostEndpoint);
    socket.onmessage = function(evt){
        // console.log(evt.data);
        let map =JSON.parse(evt.data);
        drawMap(map);
    };

    socket.onopen = () => {
        console.log("Successfully Connected");
        // socket.send("Hi From the Client!");
    };

    socket.onclose = event => {
        console.log("Socket Closed Connection: ", event);
        // socket.send("Client Closed!");
    };

    socket.onerror = error => {
        console.log("Socket Error: ", error);
    };
}

function clickEvent(event){
    let elemX = Math.min(Math.round((event.pageX - canvasOffsetLeft) / cellW), golWidth-1);
    let elemY = Math.min(Math.round((event.pageY - canvasOffsetTop) / cellH), golHeight-1);
    console.log("mpousecClick", elemX, elemY);
    socket.send(JSON.stringify({nextGen: false, x:elemX, y:elemY}));
}

function nextGenEvent(event){
    console.log("nextGenEvent");
    socket.send(JSON.stringify({nextGen: true, x:0, y:0}));
}