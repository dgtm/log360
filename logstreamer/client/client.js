const { LogRequest, LogResponse } = require("./logstreamer_pb.js");
const { LogStreamerClient } = require("./logstreamer_grpc_web_pb.js");
console.log("where is server")

var client = new LogStreamerClient('http://localhost:8080');

var request = new LogRequest();

request.setMinutes(2);

client.processRequest(request, {}, (err, response) => {
    console.log(request)
    console.log(err)

    console.log("Result of minutes : ", response.getResult())
})
// var logService = new proto.logstreamer.LogStreamerClient('http://localhost:50551');

// var request = new proto.logstreamer.LogRequest();
// request.setMinutes(20);
// var metadata = {'custom-header-1': 'value1'};
// logService.processRequest(request, metadata, function(err, response) {
//   if (err) {
//     console.log(err.code);
//     console.log(err.message);
//   } else {
//     console.log(response.getResult());
//   }
// });