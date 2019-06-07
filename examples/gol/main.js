
// function refresh(fun, millis){
//     fun();
//     setInterval(fun, millis);
// }



var hostEndpoint = "ws://localhost:8000/ws";
var socket = null;
var canvas = null;
var canvasOffsetLeft = null;
var canvasOffsetTop = null;
var golWidth = 50;
var golHeight = 50;
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
}

function clearCanvas(ctx) {
    ctx.fillStyle = backgroundStyle;
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    ctx.fillStyle = cellStyle;
}

function drawCell(ctx, x, y){
    ctx.fillRect(x * cellH, y * cellW, cellH, cellW);
}

function drawPopulation(coord){
    let ctx = canvas.getContext("2d");
    clearCanvas(ctx);
    for (let idx in coord){
        drawCell(ctx, coord[idx][0], coord[idx][1]);
    }
}

function connectToHost(){
    socket = new WebSocket(hostEndpoint);
    socket.onmessage = function(evt){
        // console.log(evt.data);
        let coord =JSON.parse(evt.data);
        drawPopulation(coord);
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
    console.log("click", elemX, elemY);
    socket.send(JSON.stringify({x:elemX, y:elemY}));
}