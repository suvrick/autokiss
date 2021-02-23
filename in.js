
var selfID = 0;
var sendData = {};
var _this = {};

//var urlData = "http://localhost:8080/autokiss/who";
//var urlInit = "http://localhost:8080/autokiss/init";
 var urlData = "https://suvricksoft.ru/autokiss/who";
 var urlInit = "https://suvricksoft.ru/autokiss/init";

function getData(data) {
    
    var xhr = new XMLHttpRequest();
    xhr.open("POST", urlData,  true);
    xhr.setRequestHeader('Content-Type', 'application/octet-stream');
    xhr.onload = function() {

        if(xhr.status === 403){
            alert("Программа не зарегистрирована.Тестовый пириод закончился.\nНапишите мне в телеграмм @help_auto_kiss для приобретения программы");
        }

        if(xhr.status === 200){
            var result = JSON.parse(xhr.responseText);
            if(result.code != 0) {
                setTimeout(callHandler.bind(_this, result), result.delay);
            }
        }
    };
    xhr.send(data);
}

function callHandler(result){

    if (_this.hasOwnProperty("Main")) {
       _this.Main.connection.sendData(result.code, result.data);
       return;
    }

    if (_this.hasOwnProperty("Game")) {
        core.protocol.Connection.sendData(result.code, result.data);
        return;
     }
}

function init() {

    if (this.hasOwnProperty("Main")) {
        _this = this;
        selfID = Main.self.id;
        Main.connection.listen(receiveDataMain, [28, 29]);
        sendData = Main.connection.sendData;
    }

    if (this.hasOwnProperty("Game")) {
        _this = this;
        selfId = Game.selfId;
        core.protocol.Connection._instance.receiveData = receiveDataGame;
        sendData = core.protocol.Connection.sendData;
    }

    var xhr = new XMLHttpRequest();
    xhr.open("GET", urlInit + "/" + selfID,  true);
    xhr.send();

}

function receiveDataMain(buffer) {

   var arr = new ArrayBuffer(buffer.bytesLength + 6);
   var data = new DataView(arr, 0, buffer.bytesLength + 6);

   data.setInt32(0, buffer.id, true);
   data.setInt16(4, buffer.type, true);

   if(buffer.type === 29 ) {
        data.setInt32(6, buffer[0], true);
        data.setInt32(10, buffer[1], true);
        data.setInt32(14, buffer[2], true);
        data.setInt32(18, buffer[3], true);
   }

   if(buffer.type === 28 ) {
        data.setInt32(6, buffer[0], true);
  }

  getData(data.buffer);

}

function receiveDataGame(buffer) {
        var id = buffer.readInt();

        var type = (buffer.readUnsignedShort()) | 0;
        if (type >= 325 || type <= 0) {
            framework.utils.dev.Logger.sendError([("Received server packet with wrong type " + type)]);
            return;
        }

        this.packetsQueue.fastPush(new core.protocol.PacketServer(id, type, buffer));
        if ([28, 29].includes(type)){
            getData(buffer.buffer_u8.slice(0, buffer._position));
        }
        this.receivePacket();
}


init();