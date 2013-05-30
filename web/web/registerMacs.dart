
import 'package:web_ui/web_ui.dart';
import 'dart:html';
import 'dart:json';
import 'dart:async';

class registerMacs extends WebComponent {
  
  
  @observable
  String cid;
  String mac;
 // List<String> macs = toObservable(new List<String>());
  
  registerMacs(){
     // new Timer.periodic(new Duration(seconds:1), (Timer t) => makeRequest());
  }
  
  void makeRequest() {
    print("makerequers"+"");
    var path = "http://127.0.0.1:8080";
    var stuff = '?"site"="registerMacAddress"';
    var json = '&"json"={"cid":"$cid","mac":"$mac"}';
    print(path+stuff+json);
    var request = HttpRequest.getString(path+stuff+json).then(requestComplete);
  }
  
  void requestComplete(String request) {
  
  }
  
}

  

