
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
    themap = "0".repeat(golWidth * golHeight);

    canvas = document.getElementById("GOLCanvas");
    let ctx = canvas.getContext("2d");

    cellW = canvas.width / golWidth;
    cellH = canvas.height / golHeight;
    canvasOffsetLeft = canvas.offsetLeft;
    canvasOffsetTop = canvas.offsetTop;

    clearCanvas(ctx);

    // canvas.addEventListener('click', clickEvent);
    canvas.addEventListener('mousemove', mouseMove);
    canvas.addEventListener('mousedown', mouseDown);
    canvas.addEventListener('mouseup', mouseUp);

    document.getElementById("nextGenButton").addEventListener('click', nextGenEvent);
    document.getElementById("clearButton").addEventListener('click', clearEvent);
    document.getElementById("randomizeButton").addEventListener('click', randomizeEvent);
    document.getElementById("randomizeGlidersButton").addEventListener('click', randomizeGlidersEvent);
}

function clearCanvas(ctx) {
    ctx.fillStyle = backgroundStyle;
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    ctx.fillStyle = cellStyle;
}

function drawCell(ctx, x, y){
    // console.log("drawCell", x, y)
    ctx.fillRect(x * cellH, y * cellW, cellH, cellW);
}

var themap = null;

function drawMap(){
    let ctx = canvas.getContext("2d");
    clearCanvas(ctx);
    // console.log("map data: "+themap);
    for (let i = 0; i < themap.length; i++){
        if (themap.charAt(i) == "1"){
            drawCell(ctx, i % golWidth, Math.floor(i / golWidth));
        }
    }
}

function updateMapAt(x, y){
    if (themap != null){
        let idx = y * golWidth + x;
        if (idx > 0 && idx < themap.length){
            let c = themap.charAt(idx) == "1" ? "0" : "1";
            themap = themap.slice(0, idx) + c + themap.slice(idx+1, themap.length);
        }
    }
}

function connectToHost(){
    socket = new WebSocket("ws://" + location.host + "/ws");
    socket.onmessage = function(evt){
        // let map =JSON.parse(evt.data);
        // drawMap(map);
        themap = String(evt.data);
        drawMap();
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

function nextGenEvent(event){
    console.log("nextGenEvent");
    socket.send(JSON.stringify({cmd: 1, x:0, y:0}));
}

function clearEvent(event){
    console.log("clearEvent");
    socket.send(JSON.stringify({cmd: 2, x:0, y:0}));
}

function randomizeEvent(event){
    console.log("randomizeEvent");
    socket.send(JSON.stringify({cmd: 3, x:0, y:0}));
}

function randomizeGlidersEvent(event){
    console.log("randomizeGlidersEvent");
    socket.send(JSON.stringify({cmd: 4, x:0, y:0}));
}

var dragState = false;
var coords = [];

function mouseDown(event){
    dragState = true;
    // console.log("mouseDown")
    coords = [];
    genCoord(event);
}

function mouseUp(event){
    // console.log("mouseUp");
    if (dragState){
        console.log("send map update to server")
        socket.send(JSON.stringify({cmd: 0, coord:coords}));
        coords = [];
    }
    dragState = false;
}

function mouseMove(event){
    // console.log("mouseMove")
    if (dragState){
        genCoord(event);
    }
}

function genCoord(event){
    let elemX = Math.floor(Math.min((event.pageX - canvas.offsetLeft) / cellW, golWidth-1)) ;
    let elemY = Math.floor(Math.min((event.pageY - canvas.offsetTop) / cellH, golHeight-1)) ;
    if (coords.length == 0){
        coords.push({x:elemX, y:elemY});
        updateMapAt(elemX, elemY);
    } else {
        let lastX = coords[coords.length-1].x;
        let lastY = coords[coords.length-1].y;
        if ( lastX!= elemX || lastY != elemY){
            coords.push({x:elemX, y:elemY});
            updateMapAt(elemX, elemY);
        }
    }
    drawMap();
}

