
var _this = {};

var isDebug = true;
var urlData = "";
var urlInit = "";

var selfID = 0;
var roomID = 0;

var isGame = false;
var isMenuShow = false;
var autoSaveKick = false;
var autoMoveToRoom = true;
var isShowAlert = false;

var menu = {}
var screenGame = {}
var headerButtons = {}
var btnClose = {}
var autoKissBtn = {}
var autoSaveBtn = {}
var hidePopupBtn = {}





function getData(data) {

    var xhr = new XMLHttpRequest();
    xhr.open("POST", urlData, true);
    xhr.setRequestHeader('Content-Type', 'application/octet-stream');
    xhr.onload = function () {

        if (xhr.status === 403) {
            showAlert();
        }

        if (xhr.status === 200) {
            var result = JSON.parse(xhr.responseText);
            if (result.code != 0) {
                setTimeout(callHandler.bind(_this, result), result.delay);
            }
        }
    };
    xhr.send(data);
}

function callHandler(result) {
    _this.Main.connection.sendData(result.code, result.data);
}

// Send packet MOVE type 21 id 259 data: 21 (2) [22132982, 0]


//packet KICK_KICKS:308 with id 1894 and length 1 data: 17770939,22132982,30
//Send packet BOTTLE_SAVE type 30 id 604 data: 30 [17770939]
function showAlert() {

    if (isShowAlert)
        return;

    var div = document.createElement("div")
    div.style.position = "absolute";
    div.style.bottom = "0"
    div.style.width = "300px"
    div.style.padding = "10px";
    div.style.background = "white";

    var p = document.createElement("span");
    p.style.padding = " 0 10px 10px 10px";
    p.style.display = "block";
    p.innerText = "Программа не зарегистрирована.Тестовый пириод закончился.\nНапишите мне в телеграмм @help_auto_kiss для приобретения программы";

    var head = document.createElement("span");
    head.innerText = "Bottle Auto Kiss Helper";
    head.style.display = "block";
    head.style.fontWeight = "bold";
    head.style.paddingLeft = "10px"

    var btnClose = document.createElement("span");
    btnClose.style.display = "block";
    btnClose.style.position = "absolute";
    btnClose.style.right = "10px";
    btnClose.style.top = "10px";
    btnClose.style.cursor = "pointer";
    btnClose.innerText = "x";
    btnClose.addEventListener("click", function () {
        root.removeChild(div)
    })

    div.appendChild(head);
    div.appendChild(p);
    div.appendChild(btnClose);

    screenGame.appendChild(div)

    isShowAlert = true;
}

function receiveDataMain(buffer) {

    if(buffer.type === 25){
        roomID = buffer[0]
        console.log("RoomID:", roomID)
        return;
    }


    if(!isGame)
    return; 

    var arr = new ArrayBuffer(buffer.bytesLength + 6);
    var data = new DataView(arr, 0, buffer.bytesLength + 6);

    data.setInt32(0, buffer.id, true);
    data.setInt16(4, buffer.type, true);

    if (buffer.type === 29) {
        data.setInt32(6, buffer[0], true);
        data.setInt32(10, buffer[1], true);
        data.setInt32(14, buffer[2], true);
        data.setInt32(18, buffer[3], true);
    }


    if(buffer.type === 27 ){
        console.log(buffer[0],selfID , autoMoveToRoom)
        if(buffer[0] === selfID && autoMoveToRoom) {
            console.log("tru move to roomID:", roomID)

            if(roomID === 0) {
                return;
            }
            
            Main.connection.sendData(202, roomID)
            return;
        }
    }

    if (buffer.type === 28) {
        data.setInt32(6, buffer[0], true);
    }

    if (buffer.type === 308) {

        if(buffer.bytesLength < 3) {
            console.log("308 >>>>>>>> ", buffer)
            return;
        }

        var kickID = buffer[0][0][0]
        if (kickID != selfID ){
            return;
        }

        if(!autoSaveKick) {
            return;
        }

        data.setInt32(6, buffer[0][0][0], true);
        data.setInt32(10, buffer[0][0][1], true);
        console.log("autosavekick send");
    }

    getData(data.buffer);

}

function setTopLine() {
    document.getElementsByTagName("body")[0].style.borderTop = "3px solid yellow";
}

function delTopMark() {
    document.getElementsByTagName("body")[0].style.borderTop = "0px solid yellow";
}

function createPopupMenu(){
    menu = document.createElement("div")
    menu.classList.add("menu")

    menu.innerHTML = `
    
    <ul>
        <li>
            <h2>Настроки</h2>
            <span id="btnClose">x</span>
        <li>
        <li>
            <label id="autoKissBtn" >Автопоцелуи (вкл)</label>
        </li>
        <li>
            <label id="autoSaveBtn" >Автоспасения (вкл)</label>
        </li>
        <li>
            <label id="hidePopupBtn" >Скрыть всплыв.окна (вкл)</label>
        </li>
    </ul>
    `

    screenGame.appendChild(menu);
}

function addBtn() {


    btnClose = document.getElementById("btnClose")
    autoKissBtn = document.getElementById("autoKissBtn")
    autoSaveBtn = document.getElementById("autoSaveBtn")
    hidePopupBtn = document.getElementById("hidePopupBtn")


     var btn = document.createElement("span")
     btn.innerText = "⋮";
     btn.classList.add("btn")
     btn.addEventListener("click", function(){
        isMenuShow = !isMenuShow;
        if(isMenuShow){
            menu.style.display = "block";
        } else {
            menu.style.display = "none";
        }

     })
 
     btn.addEventListener("mouseover", function(){
         btn.style.opacity = 1;
     })
 
     headerButtons.appendChild(btn)

     btnClose.addEventListener("click",function(){
        isMenuShow = false;
        menu.style.display = "none";
     })

    autoSaveBtn.addEventListener("click",function(){
        isMenuShow = false;
        menu.style.display = "none";
        autoSaveKick = !autoSaveKick;

        if(autoSaveKick) {
            autoSaveBtn.innerText = "Автоспасения (выкл)"
        } else {
            autoSaveBtn.innerText = "Автоспасения (вкл)"
        }
     })

     autoKissBtn.addEventListener("click",function(){
        isMenuShow = false;
        menu.style.display = "none";
        isGame = !isGame;
        if(isGame) {
            autoKissBtn.innerText = "Автопоцелуи (выкл)"
            setTopLine();
        } else {
            autoKissBtn.innerText = "Автопоцелуи (вкл)"
            delTopMark();
        }
     })

     hidePopupBtn.addEventListener("click",function(){
        isMenuShow = false;
        menu.style.display = "none";
     })
 
     setInterval(()=>{
         btn.style.opacity = 0.1;
     }, 3000)
 }

function init() {

    if (!this.hasOwnProperty("Main"))
        return;

    screenGame = document.getElementById("screen_game");
    if(screenGame === undefined || screenGame === null)
        return;

    headerButtons = document.getElementsByClassName("header-buttons")[0]
    if(headerButtons === undefined || headerButtons === null)
        return;

    clearInterval(timerInit)

    if(isDebug) {
        urlData = "http://localhost:8080/autokiss/who";
        urlInit = "http://localhost:8080/autokiss/init";
    } else {
        urlData = "https://suvricksoft.ru/autokiss/who";
        urlInit = "https://suvricksoft.ru/autokiss/init";
    }

    _this = this;
    selfID = Main.self.id;
    Main.connection.listen(receiveDataMain, [25, 28, 29, 308]);

    var xhr = new XMLHttpRequest();
    xhr.open("GET", urlInit + "/" + selfID, true);
    xhr.send();


    createPopupMenu();
    addBtn();
}

var timerInit = setInterval(init, 1000)

/*





function createMenu(root) {

    console.log("root", root)
    //var root = document.getElementById("screen_game");
    var div = document.createElement("div");
    div.classList.add("menu_wrapper");

    div.innerHTML = `
    <dl>
    <dt> Автопоцелуи </dt>
    <dd>Вкл</dd>

    <dt> Автосохранения </dt>
    <dd>Вкл</dd>
</dl>
    `

    // var li = document.createElement("li")
    // li.innerText = "вкл";
    // ul.appendChild(li)

    root.append(div)
}




function callShowMenu(e){
    if (e.code === Key) {
        inShowMenu = !inShowMenu;
        if(inShowMenu) {
            menu.style.display = "flex"
        } else {
            menu.style.display = "none"
        }
    }
}

addEventListener("keydown", callShowMenu);
var timer = setInterval(initUI, 1000);
*/