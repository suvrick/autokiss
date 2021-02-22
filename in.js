
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
                setTimeout(callHandler.bind(this, result), result.delay);
            }
        }
    };
    xhr.send(data);
}

function callHandler(result){
    core.protocol.Connection.sendData(result.code, result.data);
}

function init() {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", urlInit + "/" + Game.selfId,  true);
    xhr.send();
}

function receiveData(buffer) {
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

core.protocol.Connection._instance.receiveData = receiveData;
init();